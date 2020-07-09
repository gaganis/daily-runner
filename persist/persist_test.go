package persist

import (
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
	if err := os.Remove(fileName); err != nil {
		t.Fatal(err)
	}

	//act
	writeTimeToFile(previousTime, fileName)
	result := readTimeFromFile(fileName)
	t.Log(result)

	//assert
	if !result.Equal(previousTime) {
		t.Error("Unable to write and then read time to file")
	}
}
