package persist

import (
	"bufio"
	"daily-run-wrapper/environment"
	"fmt"
	"os"
	"time"
)

func LastRunTimeExists() bool {
	return lastReadFileExists(environment.GetLastRunFileName())
}

func WriteLastRunTime(timeRun time.Time) {
	writeTimeToFile(timeRun, environment.GetLastRunFileName())
}

func ReadLastRunTime() time.Time {
	return readTimeFromFile(environment.GetLastRunFileName())
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

	_, err = file.WriteString(timeRun.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
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
