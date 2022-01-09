package util

import (
	"fmt"
	"strconv"
)

func GetYearFromString(yearString string) (int, error) {
	year, err := strconv.Atoi(yearString)
	if err != nil || year < 1000 || year >= 10000 {
		return -1, fmt.Errorf("invalid year '%v'", yearString)
	}
	return year, nil
}
