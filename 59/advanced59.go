package main

import (
	"encoding/json"
	"fmt"
)

type Node struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Parent   *Node   `json:"parent,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

func main() {
	root := &Node{ID: 1, Name: "Root"}
	child1 := &Node{ID: 2, Name: "Child 1", Parent: root}
	child2 := &Node{ID: 3, Name: "Child 2", Parent: root}
	root.Children = []*Node{child1, child2}

	// Попытка сериализации, именно здесь мы и получим ошибку.
	data, err := json.Marshal(root)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Успех:", string(data))
	}
}
