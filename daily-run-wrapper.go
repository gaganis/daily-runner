package main

import (
	"cloud.google.com/go/civil"
	"daily-run-wrapper/environment"
	"daily-run-wrapper/persist"
	"daily-run-wrapper/run"
	"fmt"
	"github.com/nightlyone/lockfile"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

func shouldRun(lastRun time.Time, atTime time.Time, targetTime civil.Time) bool {

	//Round time to catch cases where we do not wakeup at precisely the target time
	roundedTime := atTime.Round(10 * time.Minute)
	if roundedTime.Equal(timeInstanceFromLocalTime(targetTime, atTime)) {
		return true
	}
	if lastRun.Add(24 * time.Hour).Before(atTime) {
		return true
	}
	return false
}

func timeInstanceFromLocalTime(localTime civil.Time, dateSource time.Time) time.Time {
	localDateTime := civil.DateTime{
		Date: civil.DateOf(dateSource),
		Time: localTime,
	}
	return localDateTime.In(time.Local)
}

func main() {

	initProfileFromArgs()

	logFile := setupWrapperLogger()
	defer func() {
		if e := logFile.Close(); e != nil {
			panic(e)
		}
	}()

	atTime, err := civil.ParseTime("01:00:00")
	if err != nil {
		panic(fmt.Errorf("could not parse targetTime %v", err))
	}
	runSingleProcess(atTime)

	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			panic(x)
		}
	}()
}

func initProfileFromArgs() {
	if len(os.Args) == 2 {
		profile := os.Args[1]

		if profile == "default" {
			fmt.Print("Wrong profile name argument provided 'default'. 'default' is reserved and cannot be used. " +
				"Please provide a different name. Exiting")
			os.Exit(1)
		}

		matched, _ := regexp.MatchString(`^[0-9a-zA-Z.-_]*$`, profile)
		if !matched {
			fmt.Printf("Wrong profile name argument provided: '%v'. Please provide a name containing only latin "+
				"letters, numbers and the characters '.-_'. Exiting", profile)
			os.Exit(1)
		}

		environment.SetProfile(profile)
	}
}

func setupWrapperLogger() *os.File {
	logFilePath := path.Join(environment.GetLocalAppDataDir(), "log/wrapper.log")

	if err := os.MkdirAll(path.Dir(logFilePath), 0755); err != nil {
		panic(fmt.Errorf("unable to create directories to write logfile %v, %v", logFilePath, err))
	}

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("error opening log file: %v", err))
	}

	log.SetOutput(f)
	return f
}

// This function employs a pid lockfile so that only one process of daily-run-wrapper is running at one time
func runSingleProcess(targetTime civil.Time) {
	lock, err := lockfile.New(filepath.Join(os.TempDir(), "daily-run-wrapper.lck"))
	if err != nil {
		panic(err) // handle properly please!
	}

	// Error handling is essential, as we only try to get the lock.
	if err = lock.TryLock(); err != nil {
		fmt.Print("Another process of daily-run-wrapper is already running, exiting.")
		os.Exit(1)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v", lock, err)
			panic(err) // handle properly please!
		}
	}()

	mainLoop(targetTime)
}

func mainLoop(targetTime civil.Time) {

	log.Print("Starting main loop")

	for {
		hasLastRun := persist.LastRunTimeExists()

		// Reading the time before starting the operation as an upload cat take many days to upload
		// and the start time is more indicative of the backup's freshness
		startTime := time.Now()

		log.Printf("Deciding whether to run with hasLastRun: %v", hasLastRun)
		if hasLastRun {
			log.Printf("Last run read from file: %v", persist.ReadLastRunTime())
		}

		if !hasLastRun || shouldRun(persist.ReadLastRunTime(), startTime, targetTime) {

			if run.WrappedCommand() {
				log.Printf("Writing time to file %v", startTime)
				persist.WriteLastRunTime(startTime)
			} else {
				log.Printf("Command failed to run")
			}
		}
		time.Sleep(4 * time.Minute)
	}
}
