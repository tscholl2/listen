package main

import (
	"bytes"
	"os/exec"
)

var (
	wavHeader = []byte{82, 73, 70, 70, 36, 125, 0, 0, 87, 65, 86, 69, 102, 109, 116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 128, 62, 0, 0, 0, 125, 0, 0, 2, 0, 16, 0, 100, 97, 116, 97, 0, 125, 0, 0}
)

func stt(sound []int16) string {
	f := make([]byte, len(wavHeader)+2*len(sound))
	for i, b := range wavHeader {
		f[i] = b
	}
	for i, x := range sound {
		f[len(wavHeader)+2*i] = uint8(x & 0xff)
		f[len(wavHeader)+2*i+1] = uint8(x >> 8)
	}
	//ioutil.WriteFile("x.wav", f, 0666)
	c := exec.Command("pocketsphinx_continuous",
		"-dict", "sphinx/2772.dic",
		"-lm", "sphinx/2772.lm",
		"-logfn", "/dev/null",
		"-infile", "/dev/stdin",
	)
	c.Stdin = bytes.NewReader(f)
	buf := new(bytes.Buffer)
	c.Stdout = buf
	c.Run()
	return string(buf.Bytes())
}
