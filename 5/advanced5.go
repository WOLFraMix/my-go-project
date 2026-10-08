package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Для чтения стандартного defer без проверки ошибок обычно достаточно
	defer file.Close()

	// Создаем буфер на 8 байт
	buf := make([]byte, 8)

	for {
		n, err := file.Read(buf)

		// Сначала обрабатываем то, что прочитали (даже если err != nil)
		if n > 0 {
			// Используем срез buf[:n], чтобы взять только реально прочитанные байты
			fmt.Printf("Прочитан кусок (%d байт): %q\n", n, string(buf[:n]))
		}

		// Затем проверяем, не достигли ли мы конца файла
		if err == io.EOF {
			fmt.Println("Достигнут конец файла.")
			break
		}

		// Проверяем другие возможные ошибки
		if err != nil {
			log.Fatal(err)
		}
	}
}
