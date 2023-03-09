package util

import (
	"math"
)

func LeqWithTolerance(a, b, tolerance float64) bool {
	if EqWithTolerance(a, b, tolerance) {
		return true
	} else {
		return a < b
	}
}

func EqWithTolerance(a, b, tolerance float64) bool {
	return math.Abs(a - b) < tolerance
}