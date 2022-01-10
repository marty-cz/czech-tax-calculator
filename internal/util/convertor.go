package util

import (
	"fmt"
	"strconv"
	"time"
)

func GetYearFromString(yearString string) (int, error) {
	year, err := strconv.Atoi(yearString)
	if err != nil || year < 1000 || year >= 10000 {
		return -1, fmt.Errorf("invalid year '%v'", yearString)
	}
	return year, nil
}

func GetDateThreeYearsBefore(date time.Time) time.Time {
	return time.Date(date.Year()-3, date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
}
