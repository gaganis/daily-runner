package main

import (
	"bufio"
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/nightlyone/lockfile"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
		civil.DateOf(dateSource),
		localTime,
	}
	return localDateTime.In(time.Local)
}

func lastReadFileExists(fileName string) bool {
	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)
}

func writeTimeToFile(timeRun time.Time, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("Unable to write lastRunFile: %v ", fileName))
	}

	file.WriteString(timeRun.Format(time.RFC3339))
}

func readTimeFromFile(fileName string) time.Time {
	open, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("Unable to read from %v", fileName))
	}

	scanner := bufio.NewScanner(open)
	scanner.Scan()
	scannedText := scanner.Text()
	parsedTime, err := time.Parse(time.RFC3339, scannedText)
	if err != nil {
		panic(fmt.Errorf("unable to parse date from string:  %v in file %v", scannedText, fileName))
	}

	return parsedTime
}

func main() {
	logFile := setupLogger()
	defer logFile.Close()

	atTime, err := civil.ParseTime("01:00:00")
	if err != nil {
		panic(fmt.Errorf("Could not parse targetTime %v", err))
	}
	runSingleProcess(atTime)

	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			panic(x)
		}
	}()
}

func setupLogger() *os.File {
	logFilePath := path.Join(getLocalAppDataDir(), "log/wrapper.log")

	if err := os.MkdirAll(path.Dir(logFilePath), 0755); err != nil {
		panic(fmt.Errorf("Unable to create directories to write logfile %v, %v", logFilePath, err))
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
		lastRunFileName := getLastRunFileName()
		hasLastRun := lastReadFileExists(lastRunFileName)

		startTime := time.Now()

		log.Printf("Deciding whether to run with hasLastRun: %v", hasLastRun)
		if hasLastRun {
			log.Printf("Last run read from file: %v", readTimeFromFile(lastRunFileName))
		}

		if !hasLastRun || shouldRun(readTimeFromFile(lastRunFileName), startTime, targetTime) {

			runWrappedCommand()
			log.Printf("Writing time to file %v", startTime)
			writeTimeToFile(startTime, getLastRunFileName())
		}
		time.Sleep(4 * time.Minute)
	}
}

func runWrappedCommand() {
	logFile, err := openLogFile()
	defer logFile.Close()

	log.Printf("Running process")
	command := exec.Command("duply", "zenbook_backup", "backup")

	output, err := command.CombinedOutput()
	fmt.Print(len(output))
	if err != nil {
		panic(err)
	}

	_, err = logFile.Write(output)
	if err != nil {
		panic(err)
	}
}

func openLogFile() (*os.File, error) {
	logFilePath := path.Join(getLocalAppDataDir(), "log/command-output.log")
	if err := os.MkdirAll(path.Dir(logFilePath), 0755); err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		panic(err)
	}
	return logFile, err
}

func getLastRunFileName() string {
	return path.Join(getLocalAppDataDir(), "last-run")
}

func getLocalAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homeDir, ".local/share/daily-run-wrapper/")
}
