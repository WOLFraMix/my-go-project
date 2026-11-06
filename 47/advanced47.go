package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.CreateTemp(".", "tempfile_*.log")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name()) // 2
	defer file.Close()           // 1
}
