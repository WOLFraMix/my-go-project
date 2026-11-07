package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// переименовываем
	err := os.Rename("old.txt", "new.txt")
	if err != nil {
		log.Fatal(err)
	}

	// удаляем
	err = os.Remove("file.log")
	if err != nil {
		log.Fatal(err)
	}

	// удаляем всё
	err = os.RemoveAll(".vscode")
	if err != nil {
		log.Fatal(err)
	}

	// создаём файл если его нету
	_, err = os.OpenFile("new.txt", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("file already exists")
		} else {
			log.Fatal(err)
		}
	}

	// копируем
	err = copyFile("old.txt", "new.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func copyFile(src, dst string) (err error) {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := destination.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(destination, source)
	return
}
