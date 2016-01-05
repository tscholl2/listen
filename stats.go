package main

import "math"

// returns mean, standard deviation as floats
func stats(b []int8) (float64, float64) {
	n := len(b)
	var sum, sumSq float64
	for _, x := range b {
		sum += float64(uint8(x))
		sumSq += float64(uint8(x)) * float64(uint8(x))
	}
	return sum / float64(len(b)), math.Sqrt((sumSq - sum*sum/float64(n)) / float64(n-1))
}
