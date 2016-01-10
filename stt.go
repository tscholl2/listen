package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func stt(sound []int16) string {
	fmt.Printf("length of sound for stt is %d\n", len(sound))
	f := make([]byte, 2*len(sound))
	for i, x := range sound {
		f[2*i] = uint8(x & 0xff)
		f[2*i+1] = uint8(x >> 8)
	}
	c := exec.Command("pocketsphinx_continuous",
		"-dict", "sphinx/4454.dic",
		"-lm", "sphinx/4454.lm",
		"-logfn", "/dev/null",
		"-infile", "/dev/stdin",
	)
	c.Stdin = bytes.NewReader(f)
	buf := new(bytes.Buffer)
	c.Stdout = buf
	c.Run()
	return string(buf.Bytes())
}
