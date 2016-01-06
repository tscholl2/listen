package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const witServerKey = ""

/*
$ curl -XPOST 'https://api.wit.ai/speech?v=20141022' \
 -i -L \
 -H "Authorization: Bearer $TOKEN" \
 -H "Content-Type: audio/wav" \
 --data-binary "@sample.wav"
*/
func stt(b []int8) string {
	arr := make([]byte, len(b))
	for i, x := range b {
		arr[i] = byte(x)
	}
	/*arr = append([]byte{82, 73, 70, 70, 100, 31, 0, 0, 87, 65, 86, 69, 102, 109,
	116, 32, 16, 0, 0, 0, 1, 0, 1, 0, 64, 31, 0, 0, 64, 31, 0, 0, 1, 0, 8, 0,
	100, 97, 116, 97, 64, 31, 0, 0}, arr...)*/
	req, err := http.NewRequest(
		"POST",
		"https://api.wit.ai/speech",
		bytes.NewReader(arr),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type",
		"audio/raw;encoding=signed-integer;bits=8;rate=8000;endian=little")
	req.Header.Set("Authorization", "Bearer "+witServerKey)
	req.Header.Set("Transfer-encoding", "chunked")
	fmt.Println("sending request")
	res, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("got response")

	msg := &struct {
		Text     string `json:"_text"`
		MsgID    string `json:"msg_id"`
		Outcomes []struct {
			Text       string  `json:"_text"`
			Confidence float64 `json:"confidence"`
			Entities   struct {
				Datetime []struct {
					Value struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"value"`
				} `json:"datetime"`
				Metric []struct {
					Metadata string `json:"metadata"`
					Value    string `json:"value"`
				} `json:"metric"`
			} `json:"entities"`
			Intent string `json:"intent"`
		} `json:"outcomes"`
	}{}
	a, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(a))
	json.NewDecoder(bytes.NewReader(a)).Decode(msg)
	fmt.Println(msg.Text)
	fmt.Println("--------mesage--------")
	fmt.Println(msg)
	fmt.Println("--------mesage--------")
	ioutil.WriteFile("./a.wav", arr, 0666)
	return "hi"
}
