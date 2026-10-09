package main

import "fmt"

func producer(ch chan<- int) { // chan<- означает «только для отправки»
	ch <- 10
	ch <- 20
	close(ch) // сообщаем, что больше ничего не будет
}

func main() {
	ch := make(chan int)

	go producer(ch)

	for v := range ch { // читаем, пока канал не закрыт
		fmt.Println("Получено:", v)
	}
	fmt.Println("Канал закрыт, всё прочитано")
}
