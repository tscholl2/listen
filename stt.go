package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
//witKey= ""
)

func stt(b []int8) string {
	// wav file header 44 bytes of something
	arr := []byte{82, 73, 70, 70, 164, 62, 0, 0, 87, 65, 86, 69, 102, 109, 116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 64, 31, 0, 0, 64, 31, 0, 0, 1, 0, 8, 0, 100, 97, 116, 97, 128, 62, 0, 0}
	for _, x := range b {
		arr = append(arr, byte(x))
	}
	//ioutil.WriteFile("x.wav", arr, 0666)
	//return "hi"
	req, err := http.NewRequest("POST", "https://api.wit.ai/speech", bytes.NewReader(arr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Authorization", "Bearer "+witKey)
	req.Header.Set("Transfer-encoding", "chunked")
	fmt.Println("sending request")
	res, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("got response")
		a, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(a))
		return ""
	}
	msg := &struct {
		Text string `json:"_text"`
	}{}
	json.NewDecoder(res.Body).Decode(msg)
	return msg.Text
}
