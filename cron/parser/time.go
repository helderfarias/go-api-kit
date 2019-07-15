package parser

import "time"

// This is equivalent to time.daysIn(m, year) an internal function that is not exported.
func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
