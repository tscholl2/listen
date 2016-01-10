package main

import (
	"fmt"

	"github.com/cocoonlife/goalsa"
)

const (
	recordingRate = 1500 // number of milliseconds to record each time
)

func listen(out chan<- string) {
	fmt.Println("listening...")
	samples := make(chan []int16)
	go record(samples)
	current := <-samples
	for {
		fmt.Println("gathering sample")
		current = <-samples
		mean, std := stats(current)
		fmt.Printf("computed stats: mean=%0.2f, std=%0.2f\n", mean, std)
		if std < 400 {
			fmt.Println("no word")
			continue
		}
		if start := wordStartIndex(current, mean, std); start != -1 {
			fmt.Printf("word starts at %d\n", start)
			if end := wordEndIndex(current[start:], mean, std); end != -1 {
				fmt.Printf("word starts at %d\n", start)
				out <- stt(current[start:end])
			} else {
				fmt.Println("need another samples")
				previous := current
				current = <-samples
				if end = wordEndIndex(current, mean, std); end != -1 {
					fmt.Printf("word ends at %d\n", end)
					out <- stt(append(previous[start:], current[:end]...))
				}
			}
		}
	}
}

func record(c chan<- []int16) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatS16LE, 16000, alsa.BufferParams{})
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("recording...")
		b := make([]int16, 16*recordingRate)
		dev.Read(b)
		fmt.Println("recorded...")
		c <- b
	}
}

func wordStartIndex(b []int16, mean, std float64) (start int) {
	for i, x := range b {
		if float64(x)-mean > std || float64(x)-mean < -std {
			if i > 4000 {
				return i - 4000
			}
			return 0
		}
	}
	return -1
}

func wordEndIndex(a []int16, mean, std float64) int {
	var normalSamples int
	for i, x := range a {
		if float64(x) < mean+std || float64(x) > mean-std {
			if normalSamples++; normalSamples > 4000 {
				return i
			}
		}
	}
	return -1
}
