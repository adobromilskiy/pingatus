package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/pinger"
	"github.com/adobromilskiy/pingatus/storage"
)

func init() {
	log.Printf("[INFO] App launched: %s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		fmt.Println("interrupt signal!")
		store, err := storage.GetMongoClient()
		if err != nil {
			log.Println("[ERROR] failed to get mongo client:", err)
			return
		}
		store.Close()
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Println("[ERROR] failed to load config:", err)
		return
	}

	pinger := pinger.NewPinger(cfg)
	pinger.Do(ctx)
	log.Println("[INFO] app finished.")
}
