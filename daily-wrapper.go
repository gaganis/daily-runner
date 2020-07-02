package main

import (
	"cloud.google.com/go/civil"
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
