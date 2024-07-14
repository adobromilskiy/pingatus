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
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/pinger"
	"github.com/adobromilskiy/pingatus/storage"
	"github.com/adobromilskiy/pingatus/webapi"
)

func init() {
	log.Printf("[INFO] App launched: %s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.Load()
	if err != nil {
		log.Println("[ERROR] failed to load config:", err)
		os.Exit(1)
	}

	store := storage.GetMongoClient(cfg)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		fmt.Println("interrupt signal!")
		store.Close(ctx)
		cancel()
	}()

	server := webapi.NewServer(cfg.WEBAPI, store)
	go server.Run(ctx)

	notifier, err := notifier.Get(cfg)
	if err != nil {
		log.Printf("[ERROR] error inititalize notifier: %v", err)
		return
	}

	p := pinger.NewPingatus(cfg, store, notifier)
	p.Do(ctx)

	log.Println("[INFO] app finished.")
}
