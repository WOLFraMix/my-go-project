package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Author struct {
	Firstname string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Book struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Author Author `json:"author"`
}

type Team struct {
	Members []string       `json:"members"`
	Scores  map[string]int `json:"scores"`
}

func main() {
	book := Book{
		Title: "Война и мир",
		Year:  1869,
		Author: Author{
			Firstname: "Лев",
			LastName:  "Толстой",
		},
	}

	jsonBook, err := json.Marshal(book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBook))

	team := Team{
		Members: []string{"Александр", "Иван", "Мария"},
		Scores:  map[string]int{"A": 100, "B": 90, "C": 80},
	}

	jsonTeam, err := json.Marshal(team)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonTeam))
}
