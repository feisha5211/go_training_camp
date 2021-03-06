package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func worker(stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("exit")
			return
		default:
			fmt.Println("running...")
			time.Sleep(3)
		}
	}
}

func main() {
	stop := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(stop <-chan bool) {
			defer wg.Done()
			worker(stop)
		}(stop)
	}
	fmt.Println("goroutine before num:", runtime.NumGoroutine())

	waitForSignal()
	close(stop)

	fmt.Println("stopping all job")
	wg.Wait()
	fmt.Println("goroutine after num:", runtime.NumGoroutine())
}

func waitForSignal() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	fmt.Println(<-sigs)
}
