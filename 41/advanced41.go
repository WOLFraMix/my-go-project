package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "output.txt"
	err := head(filename, 5)
	if err != nil {
		log.Fatal(err)
	}
}

func head(filename string, n int) error {
	if n <= 0 {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "close file: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(file)

	for range n {
		if !scanner.Scan() {
			break
		}
		fmt.Println(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("scan file: %w", err)
	}

	return nil
}
