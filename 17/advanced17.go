package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	go fmt.Println("Hello from GR")
	fmt.Println("Hello from main")
	time.Sleep(time.Second)
	fmt.Println(time.Since(t))
}
