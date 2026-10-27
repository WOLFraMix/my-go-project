package main

import (
	"fmt"
	"sync"
)

func main() {
	// Создаем входной канал
	input := make(chan int)

	// Горутина для генерации данных
	go func() {
		for i := range 20 {
			input <- i
		}
		close(input)
	}()

	// Промежуточный канал для Fan-Out
	processed := make(chan int)

	// Запускаем Fan-Out с 4 воркерами
	fanOut(input, processed, 4)

	// Создаем несколько каналов для Fan-In
	split1 := make(chan int)
	split2 := make(chan int)

	// Разделяем processed канал на два
	go func() {
		for n := range processed {
			split1 <- n
			split2 <- n
		}
		close(split1)
		close(split2)
	}()

	// Объединяем каналы обратно с Fan-In
	merged := fanIn(split1, split2)

	// Выводим результаты
	for result := range merged {
		fmt.Println(result)
	}
}

func fanOut(input <-chan int, output chan<- int, workerCount int) {
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for num := range input {
				// Имитация обработки
				result := num * 2
				output <- result
			}
		}()
	}

	// Закрываем output после завершения всех воркеров
	go func() {
		wg.Wait()
		close(output)
	}()
}

func fanIn(inputs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	output := make(chan int)

	// Функция для копирования значений из одного канала в output
	copy := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			output <- n
		}
	}

	wg.Add(len(inputs))
	for _, c := range inputs {
		go copy(c)
	}

	// Закрываем output после завершения всех горутин
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
