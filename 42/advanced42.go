package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	p := path.Join("a", "b", "c") // URL
	fmt.Println(p)

	f := filepath.Join("a", "b", "c") // cd
	fmt.Println(f)

	f2, err := filepath.Abs("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f2)

	f3, err := filepath.Rel("/my-go-project/42/output.txt", "/my-go-project/25/configs/config.go") // откуда куда
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f3)

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("Directory: %s\n", entry.Name())
		} else {
			fmt.Printf("File: %s\n", entry.Name())
		}
	}
}
