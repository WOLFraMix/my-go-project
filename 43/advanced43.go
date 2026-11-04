package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseDir := "/var/data"

	validInput := "gollum/treasure/my_precious.ring"
	_, err := saveReadFile(baseDir, validInput)
	fmt.Printf("Запрос: %s, ошибка: %q\n", validInput, err)

	invalidInput := "../../../etc/passwd"
	_, err = saveReadFile(baseDir, invalidInput)
	fmt.Printf("Запрос: %s, ошибка: %q\n", invalidInput, err)
}

func saveReadFile(baseDir, userInput string) ([]byte, error) {
	fullPath := filepath.Join(baseDir, userInput)

	rel, err := filepath.Rel(baseDir, userInput)
	if err != nil {
		return nil, fmt.Errorf("filepath rel: %w", err)
	}
	if strings.HasPrefix(rel, "..") {
		return nil, fmt.Errorf("attempt to escape")
	}

	return os.ReadFile(fullPath)
}
