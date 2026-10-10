package main

import (
	"io"
	"log"
	"os"
)

func main() {
	src, err := os.Open("source.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	dst, err := os.Create("dest.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}
