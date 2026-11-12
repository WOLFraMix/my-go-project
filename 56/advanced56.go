package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Event struct {
	Name        string          `json:"name"`
	PayloadType string          `json:"payload_type"`
	Payload     json.RawMessage `json:"payload"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	rawData := []byte(`{
		"name": "Событие",
		"payload_type": "user_created",
		"payload": {
			"id": 123,
			"name": "John Doe"
		}
	}`)

	var event Event
	if err := json.Unmarshal(rawData, &event); err != nil {
		log.Fatal(err)
	}

	fmt.Println(event)

	fmt.Println()

	jsonData := []byte(`{
	"id": 123,
	"price": "$12.99",
	"active":true,
	"name":"Widget",
	"discount":0.1
	}`)

	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
