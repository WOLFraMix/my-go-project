package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// DeduplicateLogs читает лог-поток из r, фильтрует идущие подряд дубликаты
// (с учетом того что время может меняться) и записывает очищенный лог в w.
func DeduplicateLogs(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	var previousLog string
	for scanner.Scan() {
		line := scanner.Text()
		if previousLog == "" {
			_, err := w.Write([]byte(line + "\n"))
			if err != nil {
				return err
			}
			previousLog = line[22:]
			continue
		}
		if previousLog == line[22:] {
			continue
		}
		_, err := w.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
		previousLog = line[22:]
	}
	err := scanner.Err()
	if err != nil {
		return err
	}
	return nil
}

// Функция main и все тесты будут от вас скрыты на сайте.
func main() {
	if !test1() || !test2() || !test3() || !test4() || !test5() || !test6() || !test7() {
		os.Exit(1)
	}
	fmt.Println("Все тесты успешно пройдены!")
}

// Пустой поток
func test1() bool {
	var in bytes.Buffer
	var out bytes.Buffer
	err := DeduplicateLogs(&in, &out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Сбой с ошибкой в test1: %v\n", err)
		return false
	}
	if out.Len() != 0 {
		fmt.Fprintf(os.Stderr, "Ожидали пустой вывод в test1, но что-то записалось: %q\n", out.String())
		return false
	}
	return true
}

// Обычные дубликаты
func test2() bool {
	inStr := "[2026-06-08 10:00:01] строка1\n[2026-06-08 10:00:02] строка1\n[2026-06-08 10:00:05] строка2\n[2026-06-08 10:00:05] строка2\n[2026-06-08 10:00:05] строка2\n[2026-06-08 10:00:07] строка3\n"
	expected := "[2026-06-08 10:00:01] строка1\n[2026-06-08 10:00:05] строка2\n[2026-06-08 10:00:07] строка3\n"
	r := strings.NewReader(inStr)
	var w bytes.Buffer
	if err := DeduplicateLogs(r, &w); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в test2: %v\n", err)
		return false
	}
	if w.String() != expected {
		fmt.Fprintf(os.Stderr, "В test2 не совпали ожидания.\nПолучили:\n%s\nОжидали:\n%s\n", w.String(), expected)
		return false
	}
	return true
}

// Дата и время меняется, месседж тот же
func test3() bool {
	inStr := "[2026-06-08 10:00:01] worker failed\n" +
		"[2026-06-08 10:00:02] worker failed\n" +
		"[2026-06-08 10:00:05] worker failed\n" +
		"[2026-06-08 10:00:10] connection ok\n"
	expected := "[2026-06-08 10:00:01] worker failed\n" +
		"[2026-06-08 10:00:10] connection ok\n"
	r := strings.NewReader(inStr)
	var w bytes.Buffer
	if err := DeduplicateLogs(r, &w); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в test3: %v\n", err)
		return false
	}
	if w.String() != expected {
		fmt.Fprintf(os.Stderr, "Ждали срезку дублей по дате и времени в test3.\nПолучили:\n%s\n", w.String())
		return false
	}
	return true
}

// Стриминг с задержкой (проверка на то, что функция не зависает в ожидании буфера)
func test4() bool {
	r, wPipe := io.Pipe()
	var out bytes.Buffer
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := DeduplicateLogs(r, &out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Тест стриминга: %v\n", err)
		}
		r.Close()
	}()

	go func() {
		parts := []string{
			"[2026-06-08 10:00:01] сообщение с задержкой\n",
			"[2026-06-08 10:00:02] сообщение с задержкой\n",
			"[2026-06-08 10:00:05] следующее сообщение\n",
		}

		for _, part := range parts {
			time.Sleep(10 * time.Millisecond)
			wPipe.Write([]byte(part))
		}
		wPipe.Close()
	}()

	wg.Wait()

	expected := "[2026-06-08 10:00:01] сообщение с задержкой\n[2026-06-08 10:00:05] следующее сообщение\n"
	if out.String() != expected {
		fmt.Fprintf(os.Stderr, "Test4 (стрим):\nПолучили:\n%q\nОжидалось:\n%q\n", out.String(), expected)
		return false
	}

	return true
}

// Тест на огромную строку
func test5() bool {
	// Создаем сообщение размером 70 КБ
	largeMessage := strings.Repeat("A", 70*1024)

	inStr := "[2026-06-08 10:00:01] " + largeMessage + "\n" +
		"[2026-06-08 10:00:02] " + largeMessage + "\n" +
		"[2026-06-08 10:00:03] short log после гиганта\n"

	expected := "[2026-06-08 10:00:01] " + largeMessage + "\n" +
		"[2026-06-08 10:00:03] short log после гиганта\n"

	r := strings.NewReader(inStr)
	var w bytes.Buffer

	if err := DeduplicateLogs(r, &w); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в test5: %v\n", err)
		return false
	}

	if w.String() != expected {
		fmt.Fprintf(os.Stderr, "В test5 не совпали ожидания при обработке больших строк.\n")
		return false
	}

	return true
}

// Проверка удаления только идущих подряд дубликатов
func test6() bool {
	inStr := "[2026-06-08 10:00:01] error A\n" +
		"[2026-06-08 10:00:02] error B\n" +
		"[2026-06-08 10:00:03] error A\n" + // Эта строка не должна удалиться, так как перед ней была строка с "error B"
		"[2026-06-08 10:00:04] error C\n"

	expected := inStr // Ожидаем, что ничего не удалится

	r := strings.NewReader(inStr)
	var w bytes.Buffer
	if err := DeduplicateLogs(r, &w); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в test6: %v\n", err)
		return false
	}
	if w.String() != expected {
		fmt.Fprintf(os.Stderr, "Провален test6. Удаляться должны только идущие подряд дубликаты.\nПолучили:\n%s\n", w.String())
		return false
	}
	return true
}

// Проверка корректного индекса префикса, строки отличаются только первой буквой сообщения
func test7() bool {
	inStr := "[2026-06-08 10:00:01] Connection error\n" +
		"[2026-06-08 10:00:02] connection error\n"

	expected := inStr // Обе строки должны остаться, они разные

	r := strings.NewReader(inStr)
	var w bytes.Buffer
	if err := DeduplicateLogs(r, &w); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в test7: %v\n", err)
		return false
	}
	if w.String() != expected {
		fmt.Fprintf(os.Stderr, "Провален test7. Проверьте правильность отрезания префикса (длина 22 символа).\nПолучили:\n%s\n", w.String())
		return false
	}
	return true
}
