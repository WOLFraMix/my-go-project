package main

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	fmt.Println(CountWords("Go go GO rust Rust", 2))
	// Ожидаемый вывод: map[go:3 rust:2]
}

// CountWords считает, сколько раз встречается каждое слово
func CountWords(text string, maxWorkers int) map[string]int {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	text = strings.ToLower(text)
	words := strings.Fields(text)

	// канал для воркеров
	ch := make(chan string)
	// канал для передачи локальных мап от воркеров
	resultsCh := make(chan map[string]int)

	var wg sync.WaitGroup

	// Запускаем воркеры
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := make(map[string]int)
			for word := range ch {
				local[word]++
			}
			// Отправляем локальную мапу в канал результатов
			resultsCh <- local
		}()
	}

	// Отправляем слова в канал
	for _, w := range words {
		ch <- w
	}
	// закрываем канал —
	// воркеры перестанут читать и выйдут из цикла
	close(ch)

	// ждём, пока все воркеры закончат
	wg.Wait()
	// закрываем канал результатов —
	// можно безопасно читать до конца
	close(resultsCh)

	// Агрегируем все локальные мапы в одну итоговую
	result := make(map[string]int)
	for local := range resultsCh {
		for k, v := range local {
			result[k] += v
		}
	}

	return result
}
