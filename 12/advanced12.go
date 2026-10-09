package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3) // три горутины
	ch := make(chan string)

	for id := 1; id <= 3; id++ {
		go func(gid int) {
			defer wg.Done() // сообщаем: я закончил
			for i := 1; i <= 10; i++ {
				time.Sleep(100 * time.Millisecond)
				ch <- fmt.Sprintf("Горутина %d: %d", gid, i)
			}
		}(id)
	}

	// Запускаем отдельную горутину, которая дождётся окончания всех рабочих горутин
	// и закроет канал — тогда range по ch завершится
	go func() {
		wg.Wait() // ждём, пока все 3 горутины сделают wg.Done()
		close(ch) // закрываем канал — это сигнал для range, что данных больше не будет
	}()

	for v := range ch { // читаем, пока канал не закрыт
		fmt.Println(v)
	}

	fmt.Println("Все горутины завершили работу")
}
