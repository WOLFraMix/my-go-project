package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func main() {
	filename := "./transactions.json"
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("ошибка чтения файла: %v\n", err)
		return
	}

	var users []User

	if err := json.Unmarshal(data, &users); err != nil {
		fmt.Printf("ошибка разбора JSON: %v\n", err)
		return
	}

	userAmount := make(map[int]int)
	for _, user := range users {
		userAmount[user.ID] += user.Amount
	}

	var maxID, maxAmount int
	for u, a := range userAmount {
		if a > maxAmount {
			maxAmount = a
		}
		if maxAmount == a {
			maxID = u
		}
	}

	fmt.Printf("Максимальное значение у id=%d, равное %d.", maxID, maxAmount)
}

type User struct {
	ID     int `json:"user_id"`
	Amount int `json:"amount"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	raw := struct {
		ID     any `json:"user_id"`
		Amount any `json:"amount"`
		*Alias
	}{Alias: (*Alias)(u)}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.Amount.(type) {
	case string:
		f, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("amount %q не является числом: %w", v, err)
		}
		u.Amount = f
	case float64:
		u.Amount = int(v)
	}

	switch i := raw.ID.(type) {
	case string:
		g, err := strconv.Atoi(i)
		if err != nil {
			return fmt.Errorf("ID %q не является числом: %w", i, err)
		}
		u.ID = g
	case float64:
		u.ID = int(i)
	}
	return nil
}
