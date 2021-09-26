package main

import (
	"context"
	"fmt"
	"log"
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

func pingHandle(w http.ResponseWriter, res *http.Request) {
	fmt.Fprintf(w, "{\"msg\": \"pong\"}")
}

func getBaseContext(net.Listener) context.Context {
	return errorCtx
}

func HandleHttp() error {
	go func() {
		<-errorCtx.Done()

		if err := httpServer.Shutdown(errorCtx); err != nil {
			fmt.Printf("[httpServer] shutdown httpServer error: %v\n", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("[httpServer] ListenAndServe error: %v", err)
	}
	return nil
}

func HandleSignal() error {
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	//signal.Notify(sc)

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

	if err := errorGroup.Wait(); err != nil {
		log.Printf("[main] exit: %v", err)
	}
}
