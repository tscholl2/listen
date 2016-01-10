package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	/*
		time cat x.wav | pocketsphinx_continuous -infile /dev/stdin -logfn /dev/null -lm sphinx/2772.lm -dict sphinx/2772.dic
	*/
	c := exec.Command("pocketsphinx_continuous",
		"-dict", "sphinx/4454.dic",
		"-lm", "sphinx/4454.lm",
		"-logfn", "/dev/null",
		"-infile", "/dev/stdin",
	)
	arr, _ := ioutil.ReadFile("x.wav")
	c.Stdin = bytes.NewReader(arr)
	buf := new(bytes.Buffer)
	c.Stdout = buf
	c.Run()
	fmt.Println(string(buf.Bytes()))
}
