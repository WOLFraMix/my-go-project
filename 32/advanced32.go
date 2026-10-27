package main

import (
	"context"
	"fmt"
	"time"
)

func infiniteGenerator(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // Гарантируем закрытие канала
		i := 1
		for {
			select {
			case <-ctx.Done():
				return // Выход при отмене контекста
			default:
				out <- i
				i++
				time.Sleep(700 * time.Millisecond)
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for num := range infiniteGenerator(ctx) {
		fmt.Println(num)
	}
	// Цикл завершится, когда генератор закроет канал
}
