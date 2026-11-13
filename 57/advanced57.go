package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	encoder := json.NewEncoder(os.Stdout)

	people := []Person{
		{"Alice", 30},
		{"Bob", 31},
	}

	for _, v := range people {
		if err := encoder.Encode(v); err != nil {
			log.Fatal(err)
		}
	}

	in := strings.NewReader(`{"name":"Pavel","age":36}
	{"name":"Stepan","age":27}`)
	decoder := json.NewDecoder(in)
	for {
		var p Person
		if err := decoder.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		log.Printf("%+v", p)
	}
}
