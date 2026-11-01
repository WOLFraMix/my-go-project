package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // гарантируем вызов cancel

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Горутина остановлена")
				return
			default:
				fmt.Println("Работаю...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	fmt.Println("Отмена контекста")
	cancel() // отправляем сигнал отмены

	time.Sleep(1 * time.Second) // даём время горутине завершиться
}
