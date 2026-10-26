package main

import (
	"fmt"
	"time"
)

func main() {
	// Создаем небуферизированный канал для передачи задач
	taskChannel := make(chan string)

	// Горутина-исполнитель
	go func() {
		for task := range taskChannel {
			// Получение задачи (блокируется, пока не будет отправки)
			fmt.Printf("Начало выполнения: %s\n", task)
			time.Sleep(1 * time.Second) // Имитация работы
			fmt.Printf("Завершено: %s\n", task)
		}
	}()

	// Отправляем 3 задачи последовательно
	tasks := []string{"задача 1", "задача 2", "задача 3"}
	for _, task := range tasks {
		fmt.Printf("Отправляю %s...\n", task)
		taskChannel <- task // Блокируется, пока исполнитель не получит задачу
		fmt.Printf("Отправлено: %s\n", task)
	}

	close(taskChannel)                 // Закрываем канал после отправки всех задач
	time.Sleep(100 * time.Millisecond) // Даем время на вывод
}
