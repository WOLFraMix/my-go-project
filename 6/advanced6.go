package main

import (
	"log"
	"os"
)

func main() {
	err := saveTargetData()
	if err != nil {
		log.Fatal(err)
	}
}

// Используем функцию с именованным возвращаемым значением err,
// чтобы перехватить ошибку Close() в defer
func saveTargetData() (err error) {
	// Создаем файл для записи (если файл был, он очистится)
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}

	// Безопасное закрываем, если была запись. Если Close решит
	// вернуть ошибку, мы запишем ее в возвращаемую переменную err.
	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Записываем срез байт
	_, err = file.Write([]byte("Hello"))
	if err != nil {
		return err
	}

	// Записываем строку (удобная обертка над Write)
	_, err = file.WriteString(" World!\n")
	if err != nil {
		return err
	}

	// Принудительно сбрасываем кэш ОС на диск (при критически важной записи)
	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}
