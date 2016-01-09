package main

import "math"

// returns mean, standard deviation as floats
func stats(b []int16) (float64, float64) {
	n := len(b)
	var sum, sumSq float64
	for _, x := range b {
		sum += float64(x)
		sumSq += float64(x) * float64(x)
	}
	return sum / float64(len(b)), math.Sqrt((sumSq - sum*sum/float64(n)) / float64(n-1))
}
