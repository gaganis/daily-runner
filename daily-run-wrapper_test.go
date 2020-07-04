package main

import (
	"cloud.google.com/go/civil"
	"os"
	"testing"
	"time"
)

func getTimeFromString(t *testing.T, timeStr string) time.Time {
	atTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Error(err)
	}
	return atTime
}

func getCivilTimeFromString(t *testing.T, timeStr string) civil.Time {
	atTime, err := civil.ParseTime(timeStr)
	if err != nil {
		t.Error(err)
	}
	return atTime
}

func Test_time_conversion(t *testing.T) {
	targetTime := getCivilTimeFromString(t, "01:00:00")
	atTime := getTimeFromString(t, "2020-06-29T01:04:05+03:00")
	correctTime := getTimeFromString(t, "2020-06-29T01:00:00+03:00")

	resultTime := timeInstanceFromLocalTime(targetTime, atTime)

	if !resultTime.Equal(correctTime) {
		t.Errorf("result should be %v, found %v", correctTime, resultTime)
	}
}

func Test_4min_after_target_is_true(t *testing.T) {
	previousTime := getTimeFromString(t, "2020-06-28T19:04:05+03:00")
	atTime := getTimeFromString(t, "2020-06-29T01:04:05+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	result := shouldRun(previousTime, atTime, targetTime)

	if !result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}

func Test_less_than_24_hours_is_false(t *testing.T) {
	previousTime := getTimeFromString(t, "2020-06-29T01:04:05+03:00")
	atTime := getTimeFromString(t, "2020-06-29T08:04:05+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	result := shouldRun(previousTime, atTime, targetTime)

	if result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}

func Test_more_than_24_hours_is_true(t *testing.T) {
	previousTime := getTimeFromString(t, "2020-06-24T01:04:05+03:00")
	atTime := getTimeFromString(t, "2020-06-29T08:04:05+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	result := shouldRun(previousTime, atTime, targetTime)

	if !result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}

func Test_missing_parent_dirs_is_false(t *testing.T) {
	fileName := "someconfig/mparmpoutsala"
	result := lastReadFileExists(fileName)

	if result {
		t.Errorf("file %v should not exist!", fileName)
	}
}

func Test_file_missing_is_false(t *testing.T) {
	fileName := "mparmpoutsala"
	result := lastReadFileExists(fileName)

	if result {
		t.Errorf("file %v should not exist!", fileName)
	}
}

func Test_existing_file_is_true(t *testing.T) {
	fileName := "test-resources/emptyFile"
	result := lastReadFileExists(fileName)

	if !result {
		t.Errorf("file %v should not exist!", fileName)
	}
}

func Test_write_time_and_then_read_from_file(t *testing.T) {
	//arrange
	fileName := "test-resources/dateFile"
	previousTime := getTimeFromString(t, "2020-06-24T01:04:05+03:00")
	os.Remove(fileName)

	//act
	writeTimeToFile(previousTime, fileName)
	result := readTimeFromFile(fileName)
	t.Log(result)

	//assert
	if !result.Equal(previousTime) {
		t.Error("Unable to write and then read time to file")
	}
}
