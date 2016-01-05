package main

import "github.com/cocoonlife/goalsa"

const sampleSize = 2000 // number of milliseconds to sample

func listen(out chan<- string) {
	samples := make(chan []int8)
	go record(samples)
	var current, previous []int8
	for {
		current = <-samples
		switch i := wordIndex(current); {
		case i == -1:
		case i < len(current)*1/8:
			out <- stt(append(previous, current...))
		case i > len(current)*7/8:
			previous = current
			current = <-samples
			out <- stt(append(previous, current...))
		default:
			out <- stt(current)
		}
		previous = current
	}
}

func record(c chan<- []int8) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatU8, 8000, alsa.BufferParams{})
	if err != nil {
		panic(err)
	}
	for {
		b := make([]int8, 8*sampleSize)
		dev.Read(b)
		c <- b
	}
}

func wordIndex(b []int8) int {
	mean, std := stats(b)
	var min, peaks int
	for i, x := range b {
		if float64(uint8(x)) > mean+1.5*std || float64(uint8(x)) < mean-1.5*std {
			peaks++
			if min == 0 {
				min = i
			}
		}
	}
	if std < 1 || peaks < 750 || peaks > 2500 {
		return -1
	}
	return min
}
