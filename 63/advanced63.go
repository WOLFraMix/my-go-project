package main

import (
	"fmt"
	"time"
)

func main() {
	go printNumbers("goroutine1")
	go printNumbers("goroutine2")
	printNumbers("main1")
}

func printNumbers(prefix string) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("%s: %d\n", prefix, i)
		time.Sleep(time.Nanosecond * 500)
	}
}
