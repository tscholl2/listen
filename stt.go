package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func stt(b []int8) string {
	fmt.Println(len(b))
	arr := make([]byte, len(b))
	for i, x := range b {
		arr[i] = byte(x)
	}
	// wav file header 44 bytes of something
	arr = append([]byte{82, 73, 70, 70, 164, 62, 0, 0, 87, 65, 86, 69, 102, 109, 116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 64, 31, 0, 0, 64, 31, 0, 0, 1, 0, 8, 0, 100, 97, 116, 97, 128, 62, 0, 0},
		arr...)
	//ioutil.WriteFile("x.wav", arr, 0666)
	return "hi"
	req, _ := http.NewRequest(
		"POST",
		"https://stream.watsonplatform.net/speech-to-text/api/v1/recognize?"+
			"keywords=computer"+
			"&keywords_threshold=0.75"+
			"&model=en-US_NarrowbandModel",
		bytes.NewReader(arr),
	)
	req.SetBasicAuth("", "")
	req.Header.Set("Content-Type",
		"audio/wav;encoding=signed-integer;bits=8;rate=8000;endian=little;channels=1")
	fmt.Println("sending request")
	res, _ := new(http.Client).Do(req)
	fmt.Println("got response")
	a, _ := ioutil.ReadAll(res.Body)
	fmt.Println("--------mesage--------")
	fmt.Println(string(a))
	fmt.Println("--------mesage--------")
	ioutil.WriteFile("x.wav", arr, 0666)
	return "hi"
}
