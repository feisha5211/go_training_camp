package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//基于errgroup实现一个http server的启动和关闭，以及linux signal信号的注册和处理，要保证能够一个退出，全部注销退出。

func main()  {
	done := make(chan error, 2)
	stop := make(chan struct{})

	var g errgroup.Group
	g.Go(func() error {
		err := serveDebug(stop)
		fmt.Println(err, "debug")
		done <- err
		return err
	})
	g.Go(func() error {
		err := serveApp(stop)
		fmt.Println(err, "app")
		done <- err
		return err
	})

	var stopped bool
	for i := 0; i < cap(done); i++ {
		fmt.Println(stopped)
		if err := <-done; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
	//等待系统信号
	waitForSignal()
	fmt.Println("stopping all")
	close(stop)
	time.Sleep(3*time.Second)
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec")
	} else {
		fmt.Println("failed,", err)
	}
}

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr: addr,
		Handler: handler,
	}
	go func() {
		fmt.Println("pending", addr)
		<-stop
		fmt.Println("server shutdown", addr)
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello World~")
	})
	for {
		select {
		case <-stop:
			return errors.New("stop_app")
		default:
			return serve("0.0.0.0:8090", mux, stop)
		}
	}
}

func serveDebug(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "debug info")
	})
	for {
		select {
		case <-stop:
		default:
			return serve("0.0.0.0:8091", mux, stop)
		}
	}
}

func waitForSignal() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	fmt.Println(<-sigs)
}