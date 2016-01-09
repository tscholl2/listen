package main

import (
	"os/exec"
	"strings"
)

var (
	wav = []byte{82, 73, 70, 70, 36, 125, 0, 0, 87, 65, 86, 69, 102, 109, 116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 128, 62, 0, 0, 0, 125, 0, 0, 2, 0, 16, 0, 100, 97, 116, 97, 0, 125, 0, 0}
)

const (
	key = "okay computer"
)

func containsPhrase(b []int16) string {
	c := exec.Command("pocketsphinx_continuous",
		"-dict sphinx/2772.dic -lm sphinx/2772.lm -infile /dev/stdin")
	w, err := c.StdinPipe()
	if err != nil {
		panic(err)
	}
	w.Write(wav)
	for _, x := range b {
		var h, l uint8 = uint8(x >> 8), uint8(x & 0xff)
		w.Write([]byte{h, l})
	}
	out, err := c.Output()
	if strings.Contains(string(out), "okay computer") {
		return key
	}
	return ""
}
