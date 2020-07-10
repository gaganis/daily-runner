package main

import (
	"cloud.google.com/go/civil"
	"daily-run-wrapper/configuration"
	"testing"
	"time"
)

func GetTimeFromString(t *testing.T, timeStr string) time.Time {
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
	atTime := GetTimeFromString(t, "2020-06-29T01:04:05+03:00")
	correctTime := GetTimeFromString(t, "2020-06-29T01:00:00+03:00")

	resultTime := timeInstanceFromLocalTime(targetTime, atTime)

	if !resultTime.Equal(correctTime) {
		t.Errorf("result should be %v, found %v", correctTime, resultTime)
	}
}

func Test_4min_after_target_is_true(t *testing.T) {
	previousTime := GetTimeFromString(t, "2020-06-28T19:04:05+03:00")
	atTime := GetTimeFromString(t, "2020-06-29T01:04:00+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	c := configuration.Configuration{
		PreferedRunTime: targetTime,
		Interval:        4 * time.Minute,
	}
	result := shouldRun(previousTime, atTime, c)

	if !result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}

func Test_less_than_24_hours_is_false(t *testing.T) {
	previousTime := GetTimeFromString(t, "2020-06-29T01:04:05+03:00")
	atTime := GetTimeFromString(t, "2020-06-29T08:04:05+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	c := configuration.Configuration{
		PreferedRunTime: targetTime,
		Interval:        4 * time.Minute,
	}
	result := shouldRun(previousTime, atTime, c)

	if result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}

func Test_more_than_24_hours_is_true(t *testing.T) {
	previousTime := GetTimeFromString(t, "2020-06-24T01:04:05+03:00")
	atTime := GetTimeFromString(t, "2020-06-29T08:04:05+03:00")
	targetTime := getCivilTimeFromString(t, "01:00:00")

	c := configuration.Configuration{
		PreferedRunTime: targetTime,
		Interval:        4 * time.Minute,
	}
	result := shouldRun(previousTime, atTime, c)

	if !result {
		t.Errorf("should run when less than 24 hours %v, %v, %v", atTime, targetTime, previousTime)
	}
}
