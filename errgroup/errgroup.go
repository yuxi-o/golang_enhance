package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var errorGroup *errgroup.Group
var errorCtx context.Context
var httpServer *http.Server
var signalChan chan os.Signal = make(chan os.Signal, 10)
var shutdownChan chan struct{} = make(chan struct{})

func pingHandle(w http.ResponseWriter, res *http.Request) {
	fmt.Fprintf(w, "{\"msg\": \"pong\"}")
}
func shutdownHandle(w http.ResponseWriter, res *http.Request) {
	shutdownChan <- struct{}{}
	close(shutdownChan)
}

func getBaseContext(net.Listener) context.Context {
	return errorCtx
}

func HandleHttp() error {
	go func() {
		select {
		case <-errorCtx.Done():
		case <-shutdownChan:
			break
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(timeoutCtx); err != nil {
			fmt.Printf("[httpServer] shutdown httpServer error: %v\n", err)
		}
	}()

	return httpServer.ListenAndServe()
}

func HandleSignal() error {
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

LOOP:
	for {
		select {
		case <-errorCtx.Done():
			break LOOP
		case sg := <-signalChan:
			switch sg {
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				return fmt.Errorf("[signal] exit from signal(%v)", sg)
			case syscall.SIGUSR1:
				fmt.Printf("[signal] deal with signal(%v)\n", sg)
			default:
				fmt.Printf("[signal] Unknown signal(%v)\n", sg)
			}
		}
	}

	return nil
}

func main() {
	errorGroup, errorCtx = errgroup.WithContext(context.Background())

	httpServeMux := http.NewServeMux()
	httpServeMux.HandleFunc("/ping", pingHandle)
	httpServeMux.HandleFunc("/shutdown", shutdownHandle)

	httpServer = &http.Server{
		Addr:           ":8088",
		Handler:        httpServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    getBaseContext,
	}

	errorGroup.Go(HandleSignal)
	errorGroup.Go(HandleHttp)

	fmt.Printf("[main] exit: %v", errorGroup.Wait())
}
