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
	//var current, previous []int8
	var s []int8
	for {
		fmt.Println("gathering sample")
		s = <-samples
		if wordIndex(s) != -1 {
			out <- stt(s)
		}
		/*
			current = <-samples
			switch i := wordIndex(current); {
			case i == -1:
				fmt.Println("no word")
			case i < len(current)*1/4:
				fmt.Println("early word")
				out <- stt(append(previous, current...))
			case i > len(current)*3/4:
				fmt.Println("late word")
				previous = current
				current = <-samples
				out <- stt(append(previous, current...))
			default:
				fmt.Println("mid word")
				out <- stt(current)
			}
			previous = current
		*/
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
