package main

import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/nightlyone/lockfile"
	"log"
	"os"
	"path"
	"time"
)

func shouldRun(lastRun time.Time, atTime time.Time, configuration Configuration, loc *time.Location) bool {
	if configuration.HasPreferredRunTime {
		preferredRunTimeInstance := timeInstanceFromLocalTime(configuration.PreferredRunTime, atTime, loc)
		upperLimitTime := preferredRunTimeInstance.Add(configuration.Interval + 10*time.Second)
		if preferredRunTimeInstance.Before(atTime) &&
			upperLimitTime.After(atTime) {
			return true
		}
	}

	if lastRun.Add(24 * time.Hour).Before(atTime) {
		return true
	}
	return false
}

func timeInstanceFromLocalTime(localTime civil.Time, dateSource time.Time, loc *time.Location) time.Time {
	localDateTime := civil.DateTime{
		Date: civil.DateOf(dateSource),
		Time: localTime,
	}
	return localDateTime.In(loc)
}

func main() {

	configuration := ParseConfigFromFlags()

	SetProfile(configuration.Profile)
	fmt.Printf("Starting daily-run-wrapper with configuration:\n%+v\n", configuration)
	fmt.Println("Please see logs at: ", WrapperLogFilePath())
	logFile := setupWrapperLogger()
	defer func() {
		if e := logFile.Close(); e != nil {
			panic(e)
		}
	}()
	log.Printf("Starting daily-run-wrapper with configuration:\n%+v\n", configuration)

	runSingleProcess(configuration)

	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			panic(x)
		}
	}()
}

func setupWrapperLogger() *os.File {
	logFilePath := WrapperLogFilePath()

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
func runSingleProcess(configuration Configuration) {
	lock, err := lockfile.New(LockFilePath())
	if err != nil {
		panic(err) // handle properly please!
	}

	// Error handling is essential, as we only try to get the lock.
	if err = lock.TryLock(); err != nil {
		fmt.Println("Another process of daily-run-wrapper for profile " + GetProfile() + " is already running.")
		os.Exit(1)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v", lock, err)
			panic(err)
		}
	}()

	mainLoop(configuration)
}

func mainLoop(configuration Configuration) {

	log.Print("Starting main loop")

	for {
		hasLastRun := LastRunTimeExists()

		// Reading the time before starting the operation as an upload cat take many days to upload
		// and the start time is more indicative of the backup's freshness
		startTime := time.Now()

		log.Printf("Deciding whether to run with hasLastRun: %v", hasLastRun)
		if hasLastRun {
			log.Printf("Last run read from file: %v", ReadLastRunTime())
		}

		if !hasLastRun || shouldRun(ReadLastRunTime(), startTime, configuration, time.Local) {

			if WrappedCommand(configuration) {
				log.Printf("Writing time to file %v", startTime)
				WriteLastRunTime(startTime)
			} else {
				log.Printf("Command failed to run")
			}
		}
		time.Sleep(configuration.Interval)
	}
}
