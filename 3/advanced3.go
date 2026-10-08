package main

import (
	"bytes"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("grep", "go")

	// Создаем и записываем данные в STDIN команды
	var stdin bytes.Buffer
	cmd.Stdin = &stdin
	stdin.Write([]byte("golang\npython\njava\n"))

	// Буферы для STDOUT и STDERR
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Запускаем команду
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка выполнения: %s\nStderr: %s", err, stderr.String())
	}

	log.Printf("Stdout: %s", stdout.String())
}
