package main

import (
	"fmt"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // когда функция завершится, уменьшаем счётчик на 1
	for i := 1; i <= 3; i++ {
		fmt.Printf("Горутина %d: шаг %d\n", id, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2) // мы запускаем 2 горутины

	go worker(1, &wg)
	go worker(2, &wg)

	wg.Wait() // ждём, пока обе горутины сделают wg.Done()
	fmt.Println("Все горутины завершены")
}
