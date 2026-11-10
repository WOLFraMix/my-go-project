package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Operation struct {
	Name    string  `json:"name"`
	Success Success `json:"success"`
}

type Success bool

func (s Success) MarshalJSON() ([]byte, error) {
	if s {
		return []byte(`"✓"`), nil
	}
	return []byte(`"✗"`), nil
}

func (s *Success) UnmarshalJSON(data []byte) error {
	// В случае null оставляем значение false
	if string(data) == "null" {
		*s = false
		return nil
	}
	// Параметр data - это весь JSON-токен, включая кавычки
	if string(data) == `"✓"` {
		*s = true
		return nil
	}
	if string(data) == `"✗"` {
		*s = false
		return nil
	}
	return fmt.Errorf("неизвестное значение: %s", data)
}

func main() {
	operation := Operation{
		Name:    "Купить корм для кота",
		Success: true,
	}

	data, err := json.Marshal(operation)
	if err != nil {
		log.Fatalf("marshal: %s", err)
	}
	fmt.Printf("Запакованный объект: %s\n", data)

	jsonData := `{"name":"Купить корм для кота","success":"✓"}`
	var newOperation Operation
	if err := json.Unmarshal([]byte(jsonData), &newOperation); err != nil {
		log.Fatalf("unmarshal: %s", err)
	}
	fmt.Printf("Распакованный объект: %+v\n", newOperation)
}
