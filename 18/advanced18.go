package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	t := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			getHttpCode()
		}()

	}
	wg.Wait()
	fmt.Println(time.Since(t))
}

func getHttpCode() {
	resp, err := http.Get("https://google.com")
	if err != nil {
		fmt.Printf("Ошибка %s\n", err.Error())
	}
	fmt.Printf("Код: %d\n", resp.StatusCode)
}
