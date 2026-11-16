package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	var obj Person
	json.NewDecoder(os.Stdin).Decode(&obj)

	fmt.Println(obj.Firstname)
}

// Get-Content input.txt | go run main.go

// go run main.go --from-file input.txt

// go run main.go < input.txt
