package stats

import (
	"fmt"
	"time"
)

// GetCurrentPeriod returns the current month and week identifiers
func GetCurrentPeriod() (month string, week string) {
	now := time.Now()
	return GetMonthPeriod(now), GetWeekPeriod(now)
}

// GetMonthPeriod returns the month period identifier (YYYY-MM)
func GetMonthPeriod(t time.Time) string {
	return t.Format("2006-01")
}

// GetWeekPeriod returns the ISO week period identifier (YYYY-WNN)
func GetWeekPeriod(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// IsSamePeriod checks if two times are in the same period
func IsSamePeriod(t1, t2 time.Time, periodType string) bool {
	switch periodType {
	case "month":
		return GetMonthPeriod(t1) == GetMonthPeriod(t2)
	case "week":
		return GetWeekPeriod(t1) == GetWeekPeriod(t2)
	default:
		return false
	}
}

// ShouldResetPeriod checks if a score's period has expired
func ShouldResetPeriod(scorePeriod string, periodType string) bool {
	now := time.Now()
	var currentPeriod string

	switch periodType {
	case "month":
		currentPeriod = GetMonthPeriod(now)
	case "week":
		currentPeriod = GetWeekPeriod(now)
	default:
		return false
	}

	return scorePeriod != currentPeriod
}
