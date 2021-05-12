package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("goroutine before num:", runtime.NumGoroutine())
	stop := make(chan bool)
	var wg sync.WaitGroup
	for i:=0; i<2; i++ {
		wg.Add(1)
		go func(stop chan bool, i int) {
			defer wg.Done()
			startServer(stop, i)
		}(stop, i)
	}
	fmt.Println("goroutine before num:", runtime.NumGoroutine())
	waitSign()
	close(stop)
	time.Sleep(5 * time.Second)
	wg.Wait()
	fmt.Println("goroutine after num:", runtime.NumGoroutine())
	fmt.Println("stopping all ...")
}

func startServer(stop chan bool, i int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "debug info")
	})
	addr := fmt.Sprintf("0.0.0.0:809%d", i)
	s := http.Server{
		Addr: addr,
		Handler: mux,
	}
	for {
		fmt.Println("goroutine num:", runtime.NumGoroutine())
		select {
		case <-stop:
			fmt.Printf("stop...%d\n", i)
			s.Shutdown(nil)
			return
		default:
			fmt.Println("pending", addr)
			s.ListenAndServe()
		}
	}
}

func waitSign() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)
	fmt.Println(<-sign)
}
