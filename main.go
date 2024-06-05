package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/adobromilskiy/pingatus/pinger"
)

func init() {
	fmt.Printf("App launched.\nGOOS: %s, GOARCH: %s, GOVERSION: %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		fmt.Println("interrupt signal")
		cancel()
	}()

	pinger := pinger.NewHttpPinger("https://twst.dev", 200, time.Second, time.Second*2)
	pinger.Do(ctx)
	fmt.Println("App finished.")
}
