package main

import (
	"fmt"
	"io/ioutil"

	"github.com/cocoonlife/goalsa"
)

// for testing and reference
func writeAudio(b []int16) {
	a := []byte{82, 73, 70, 70, 36, 125, 0, 0, 87, 65, 86, 69, 102, 109, 116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 128, 62, 0, 0, 0, 125, 0, 0, 2, 0, 16, 0, 100, 97, 116, 97, 0, 125, 0, 0}
	for _, x := range b {
		a = append(a, uint8(x&0xff), uint8(x>>8))
	}
	ioutil.WriteFile("x.wav", a, 0666)
}

const sampleSize = 3000 // number of milliseconds to sample

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
			out <- stt(append(previous[sampleSize/2:], current[:sampleSize/2]...))
		default:
			fmt.Println("mid word")
			out <- stt(current)
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
	fmt.Printf("std=%02f, mean=%02f, peaks=%d\n", std, mean, peaks)
	if std < 300 {
		return -1
	}
	return min
}
