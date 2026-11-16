package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type legacyUser struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	var m map[string]legacyUser
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Println("Битый JSON - не бейте его :)")
		return
	}

	people := make([]Person, 0, len(m))
	for key, u := range m {
		idStr := strings.TrimPrefix(key, "user_")
		id, convErr := strconv.Atoi(idStr)
		if convErr != nil {
			fmt.Printf("Потерялся ID в ключе %q\n", key)
			return
		}
		people = append(people, Person{ID: id, Name: u.Name, Role: u.Role})
	}

	sort.Slice(people, func(i, j int) bool {
		return people[i].ID < people[j].ID
	})

	out, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

// Get-Content input.txt | go run advanced62.go
