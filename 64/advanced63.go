package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	go func() {
		// Плотный цикл без вызовов функций.
		// До Go 1.14 он бы захватил поток навсегда.
		// Сейчас Sysmon прервет его через 10мс.
		for range 10_000_000_000 {
			// просто крутим счетчик
		}
		fmt.Println("Долгая горутина завершена")
	}()

	time.Sleep(100 * time.Millisecond)

	// Благодаря принудительному вытеснению эта горутина выполнится вовремя.
	go func() {
		fmt.Println("Короткая горутина выполнена!")
	}()

	time.Sleep(time.Second)
}
