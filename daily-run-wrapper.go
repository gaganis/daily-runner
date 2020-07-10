package main

import (
	"cloud.google.com/go/civil"
	. "daily-run-wrapper/configuration"
	"daily-run-wrapper/environment"
	"daily-run-wrapper/persist"
	"daily-run-wrapper/run"
	"fmt"
	"github.com/nightlyone/lockfile"
	"log"
	"os"
	"path"
	"time"
)

func shouldRun(lastRun time.Time, atTime time.Time, configuration Configuration) bool {
	if configuration.HasPreferredRunTime {
		preferedRunTimeInstance := timeInstanceFromLocalTime(configuration.PreferredRunTime, atTime)
		upperLimitTime := preferedRunTimeInstance.Add(configuration.Interval + 10*time.Second)
		if preferedRunTimeInstance.Before(atTime) &&
			upperLimitTime.After(atTime) {
			return true
		}
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

	configuration := ParseConfigFromFlags()

	environment.SetProfile(configuration.Profile)
	logFile := setupWrapperLogger()
	defer func() {
		if e := logFile.Close(); e != nil {
			panic(e)
		}
	}()

	runSingleProcess(configuration)

	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			panic(x)
		}
	}()
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
func runSingleProcess(configuration Configuration) {
	lock, err := lockfile.New(environment.LockFilePath())
	if err != nil {
		panic(err) // handle properly please!
	}

	// Error handling is essential, as we only try to get the lock.
	if err = lock.TryLock(); err != nil {
		fmt.Println("Another process of daily-run-wrapper for profile " + environment.GetProfile() + " is already running.")
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
		hasLastRun := persist.LastRunTimeExists()

		// Reading the time before starting the operation as an upload cat take many days to upload
		// and the start time is more indicative of the backup's freshness
		startTime := time.Now()

		log.Printf("Deciding whether to run with hasLastRun: %v", hasLastRun)
		if hasLastRun {
			log.Printf("Last run read from file: %v", persist.ReadLastRunTime())
		}

		if !hasLastRun || shouldRun(persist.ReadLastRunTime(), startTime, configuration) {

			if run.WrappedCommand(configuration) {
				log.Printf("Writing time to file %v", startTime)
				persist.WriteLastRunTime(startTime)
			} else {
				log.Printf("Command failed to run")
			}
		}
		time.Sleep(configuration.Interval)
	}
}
