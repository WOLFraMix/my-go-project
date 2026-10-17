package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	path := flag.String("file", "url.txt", "path to URL file")
	flag.Parse()
	file, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	urlSlice := strings.Split(string(file), "\n")
	respCh := make(chan int)
	errChan := make(chan error)

	for _, url := range urlSlice {
		go ping(url, respCh, errChan)
	}

	for range urlSlice {
		select {
		case err := <-errChan:
			fmt.Println(err)
		case r := <-respCh:
			fmt.Println(r)
		}
	}
}

func ping(url string, respCh chan int, errCh chan error) {
	resp, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	respCh <- resp.StatusCode
}
