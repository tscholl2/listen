package main

import (
	"fmt"

	"github.com/cocoonlife/goalsa"
)

const sampleSize = 2000 // number of milliseconds to sample

func listen(out chan<- string) {
	fmt.Println("listening...")
	samples := make(chan []int8)
	go record(samples)
	var current, previous []int8
	for {
		fmt.Println("gathering sample")
		current = <-samples
		switch i := wordIndex(current); {
		case i == -1:
			fmt.Println("no word")
		case i > len(current)/2:
			fmt.Println("late word")
			previous = current
			current = <-samples
			out <- stt(append(previous[sampleSize/2:], current[:sampleSize/2]...))
		default:
			fmt.Println("mid word")
			out <- stt(current)
		}
		previous = current
		fmt.Println("finished this sample")
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
