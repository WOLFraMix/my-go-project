package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	grCount := 5
	ch := make(chan int, grCount)

	partSize := len(arr) / grCount

	for i := 0; i < grCount; i++ {
		start := i * partSize
		end := start + partSize
		go sumParts(arr[start:end], ch)
	}

	totalSum := 0
	for range grCount {
		totalSum += <-ch
	}
	fmt.Println(totalSum)
}

func sumParts(arr []int, ch chan int) {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	ch <- sum
}
