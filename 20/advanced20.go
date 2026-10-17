package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	c := make(chan int)
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			getHttpCode(c)
		})
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	for r := range c {
		fmt.Printf("Код: %d\n", r)
	}
}

func getHttpCode(c chan int) {
	resp, err := http.Get("https://google.com")
	if err != nil {
		log.Printf("Ошибка %s\n", err)
	}
	c <- resp.StatusCode
}
