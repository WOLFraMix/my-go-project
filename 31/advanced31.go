package main

import (
	"fmt"
	"sort"
)

func main() {
	data := []int{34, 7, 23, 32, 5, 62, 31, 25, 12, 8} // Входной массив с неотсортированными числами
	sorted := parallelSort(data)
	fmt.Println("Отсортированный массив:", sorted)
}

func parallelSort(arr []int) []int {
	// Делим массив на 2 части
	part1, part2 := arr[:len(arr)/2], arr[len(arr)/2:]

	// Создаём два небуферизированных канала для получения результатов
	ch1, ch2 := make(chan []int), make(chan []int)

	// Запускаем сортировку в горутинах
	go func() { sort.Ints(part1); ch1 <- part1 }()
	go func() { sort.Ints(part2); ch2 <- part2 }()

	// Ждём результаты через select
	var sorted1, sorted2 []int
	for i := 0; i < 2; i++ {
		select {
		case sorted1 = <-ch1:
		case sorted2 = <-ch2:
		}
	}

	// Объединяем результаты
	return merge(sorted1, sorted2)
}

func merge(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	return append(append(result, a[i:]...), b[j:]...)
}
