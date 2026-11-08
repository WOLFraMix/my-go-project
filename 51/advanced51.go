package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Person struct {
	Name     string
	Age      int
	password string
}

func main() {
	person := Person{
		Name:     "John",
		Age:      25,
		password: "111",
	}

	data, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	data = []byte(`{"name": "Johny", "age": 26, "password": "123"}`)
	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}

type User struct {
	unexported string
	FirstName  string    `json:"first_name"`
	Password   string    `json:"-"`
	ID         int       `json:"id,string"`
	Age        int       `json:"age,omitempty`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
}
