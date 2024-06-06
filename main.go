package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/adobromilskiy/pingatus/config"
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

	// _, err := database.GetMongoClient(ctx, "mongodb://localhost:27017", false)
	// if err != nil {
	// 	fmt.Println("Failed to connect to MongoDB:", err)
	// 	os.Exit(1)
	// }

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Failed to load config:", err)
	}

	pinger := pinger.NewPinger(cfg)
	pinger.Do(ctx)
	fmt.Println("App finished.")
}
