package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

func main() {
	var people []Person
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		fmt.Println("Ошибка чтения строки")
		return
	}
	line := scanner.Text()

	if err := json.Unmarshal([]byte(line), &people); err != nil {
		fmt.Printf("Ошибка при декодировании JSON: %v\n", err)
		return
	}

	if !scanner.Scan() {
		fmt.Println("Ошибка чтения строки")
		return
	}
	trueCity := scanner.Text()

	for _, p := range people {
		if p.City == trueCity {
			fmt.Println(p.Name)
		}
	}
}

// Get-Content input.txt | go run advanced61.go
