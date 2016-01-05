package main

import "fmt"

func main() {
	var c chan string
	go listen(c)
	for {
		fmt.Println(<-c)
	}
}
