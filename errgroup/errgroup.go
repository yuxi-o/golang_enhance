package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func HandleSignal() error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for sg := range sc {
		switch sg {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return fmt.Errorf("[signal] exit from signal(%v)", sg)
		case syscall.SIGUSR1:
			fmt.Println("[signal] deal with signal ", sg)
		case syscall.SIGUSR2:
			fmt.Println("[signal] deal with signal ", sg)
		default:
			fmt.Println("[signal] Unknown signal", sg)
		}
	}

	return nil
}

func HandleHttp() error {
	http.HandleFunc("/ping", func(w http.ResponseWriter, res *http.Request) {
		fmt.Fprintf(w, "{\"msg\": \"pong\"}")
	})

	if err := http.ListenAndServe(":9090", nil); err != nil {
		return err
	}
	return nil
}

func main() {
	//	eg, ctx := errgroup.WithContext(context.Background())
	eg := new(errgroup.Group)

	eg.Go(HandleSignal)

	eg.Go(HandleHttp)

	if err := eg.Wait(); err != nil {
		log.Printf("[main] exit: %v", err)
	}
}
