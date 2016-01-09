package main

import (
	"fmt"

	"github.com/cocoonlife/goalsa"
)

const sampleSize = 2000 // number of milliseconds to sample

func listen(out chan<- string) {
	fmt.Println("listening...")
	samples := make(chan []int16)
	go record(samples)
	var current, previous []int16
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
			out <- containsPhrase(append(previous[sampleSize/2:], current[:sampleSize/2]...))
		default:
			fmt.Println("mid word")
			out <- containsPhrase(current)
		}
		previous = current
		fmt.Println("finished this sample")
	}
}

func record(c chan<- []int16) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatS16LE, 16000, alsa.BufferParams{})
	if err != nil {
		panic(err)
	}
	for {
		b := make([]int16, 16*sampleSize)
		dev.Read(b)
		c <- b
	}
}

func wordIndex(b []int16) int {
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
	if std < 1 || peaks < 1500 || peaks > 5000 {
		return -1
	}
	return min
}
