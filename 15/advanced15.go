package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(CompressNumbers([]int{1, 2, 3, 5, 6, 8}))
}

func CompressNumbers(nums []int) []string {
	var result []string
	if len(nums) == 0 {
		return []string{}
	}

	start := nums[0]
	end := start

	for i := 1; i < len(nums); i++ {
		if nums[i] == end+1 {
			end = nums[i]
		} else {
			// Завершаем текущий диапазон
			if start == end {
				result = append(result, strconv.Itoa(start))
			} else {
				result = append(result, fmt.Sprintf("%d-%d", start, end))
			}
			// Начинаем новый
			start = nums[i]
			end = start
		}
	}

	// Добавляем последний
	if start == end {
		result = append(result, strconv.Itoa(start))
	} else {
		result = append(result, fmt.Sprintf("%d-%d", start, end))
	}

	return result
}
