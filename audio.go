package main

import (
	"fmt"

	"github.com/cocoonlife/goalsa"
)

const sampleSize = 2000 // number of milliseconds to sample

func listen(out chan string) {
	var samples chan []int8
	var current, previous []int8
	go record(samples)
	for {
		fmt.Println("starting listen loop")
		current = <-samples
		switch i := wordIndex(current); {
		case i == -1:
			fmt.Println("new sample has no word")
		case i < len(current)*1/8:
			fmt.Println("new sample starts with word")
			out <- stt(append(previous, current...))
		case i > len(current)*7/8:
			fmt.Println("new sample ends with word")
			previous = current
			current = <-samples
			out <- stt(append(previous, current...))
		default:
			fmt.Println("new sample middle with word")
			out <- stt(current)
		}
		previous = current
		fmt.Println("ending listen loop")
	}
}

func record(c chan []int8) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatU8, 8000, alsa.BufferParams{})
	if err != nil {
		panic(err)
	}
	for {
		fmt.Printf("recording for %d seconds...\n", sampleSize/1000)
		b := make([]int8, 8*sampleSize)
		dev.Read(b)
		fmt.Println("done recording. now sending...")
		c <- b
		fmt.Println("done sending.")
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
	if peaks < 500 || peaks > 2500 {
		return -1
	}
	fmt.Printf("there might be a word starting at %d from %d\n", min, len(b))
	return min
}
