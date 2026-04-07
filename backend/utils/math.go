package utils

import "math"

// Clamp restricts a value to the 0–100 range.
func Clamp(val float64) float64 {

	if val > 100 {
		return 100
	}
	if val < 0 {
		return 0
	}
	return val

}

// Round2 rounds a float to 2 decimal places.
func Round2(val float64) float64 {

	return math.Round(val*100) / 100

}
