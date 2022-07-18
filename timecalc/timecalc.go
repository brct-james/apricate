// Package timecalc provides helper functions for working with timestamps
package timecalc

import (
	"time"

	"apricate/log"
)

// Add seconds to timestamp
func AddSecondsToTimestamp(startTime time.Time, seconds int) (time.Time) {
	duration := time.Second * time.Duration(seconds)
	log.Debug.Printf("StartTime: %v, EndTime: %v", startTime, startTime.Add(duration))
	return startTime.Add(duration)
}

// Add minutes to timestamp
func AddMinutesToTimestamp(startTime time.Time, minutes int) (time.Time) {
	duration := time.Minute * time.Duration(minutes)
	log.Debug.Printf("StartTime: %v, EndTime: %v", startTime, startTime.Add(duration))
	return startTime.Add(duration)
}