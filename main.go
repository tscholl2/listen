package main

import "fmt"

func main() {
	c := make(chan string)
	go listen(c)
	for {
		fmt.Println(<-c)
	}
}
