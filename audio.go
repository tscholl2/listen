package main

import (
	"fmt"
	"math"

	"github.com/cocoonlife/goalsa"
)

const (
	recordingRate = 1500 // number of milliseconds to record each time
	stdCutOff     = 400  // standard deviation cutoff for word testing
)

func listen(out chan<- string) {
	fmt.Println("listening...")
	samples := make(chan []int16)
	go record(samples)
	for {
		fmt.Println("gathering sample")
		a := <-samples
		mean, std := stats(a)
		fmt.Printf("computed stats: mean=%0.2f, std=%0.2f\n", mean, std)
		if std < stdCutOff {
			fmt.Println("no word")
			continue
		}
		if start := wordStartIndex(a, mean, std); start != -1 {
			fmt.Printf("word starts at %d\n", start)
			out <- stt(append(a[int(math.Max(float64(start)-4000, 0)):], <-samples...))
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
		a := make([]int16, 16*recordingRate)
		dev.Read(a)
		fmt.Println("recorded...")
		c <- a
	}
}

func wordStartIndex(a []int16, mean, std float64) (start int) {
	for i, x := range a {
		if float64(x)-mean > 1.5*std || float64(x)-mean < -1.5*std {
			if i > 8000 {
				i -= 8000
			}
			return i
		}
	}
	return -1
}
