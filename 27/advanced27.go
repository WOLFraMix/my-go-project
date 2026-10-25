package main

import (
	"io"
	"log"
	"net"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	// Раз в секунду печатаем, сколько горутин и сколько потоков ОС
	go func() {
		p := pprof.Lookup("threadcreate")
		for range time.Tick(time.Second) {
			log.Printf("goroutines=%d  OS threads=%d  GOMAXPROCS=%d",
				runtime.NumGoroutine(), p.Count(), runtime.GOMAXPROCS(0))
		}
	}()

	for {
		conn, err := ln.Accept() // ← тоже паркуется в netpoller, пока нет клиентов
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer c.Close()
	// io.Copy внутри — это c.Read() + c.Write():
	// оба вызова при EAGAIN паркуют ЭТУ горутину, освобождая поток ОС
	io.Copy(c, c)
}
