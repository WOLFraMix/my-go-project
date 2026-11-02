package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "output.txt"

	err := writeBufferedData(filename)
	if err != nil {
		log.Fatalf("Критическое значение: %v", err)
	} else {
		log.Println("Данные успешно записаны")
	}
}

func writeBufferedData(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString("Hello, world!\n")
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	err = writer.Flush() // запись
	if err != nil {
		return fmt.Errorf("flush writer: %w", err)
	}

	return nil
}
