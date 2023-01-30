package main

import (
	"errors"
	"math"
)

func KelvinToAPI(k int) (int, error) {
	if k < 2900 || k > 7000 {
		return 0, errors.New("incorrect value")
	}

	result := math.Round(1_000_000 / float64(k))
	result = math.Min(result, 344)
	result = math.Max(result, 143)

	return int(result), nil
}

func APIToKelvin(x int) int {
	k := 1_000_000 / float64(x)
	result := math.Round(k/50) * 50 // Round to nearest 50
	return int(result)
}
