package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/adobromilskiy/pingatus/internal/config"
	"github.com/adobromilskiy/pingatus/internal/notifier"
	"github.com/adobromilskiy/pingatus/internal/pinger"
	"github.com/adobromilskiy/pingatus/internal/server"
	"github.com/adobromilskiy/pingatus/internal/storage/sqlite"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, os.Kill)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	logLevel := parseLogLevel(cfg.Logger.Level)

	lg := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	if cfg.Logger.IsJSON {
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	lg.Log(ctx, slog.LevelInfo, "application started")

	ntfr, err := notifier.New(lg, cfg.Notifier) //nolint:contextcheck
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	db := sqlite.New(lg, cfg.DBDSN)

	if err = db.Open(ctx); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	defer db.Close(ctx)

	endpoint, err := sqlite.NewEndpoint(ctx, db.Get())
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	p := pinger.NewPingatus(lg, cfg.Endpoints, endpoint, ntfr)

	s := server.New(lg, endpoint, cfg.ListenAddr)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()

		lg.Log(ctx, slog.LevelInfo,
			"Received terminated signal",
		)

		return nil
	})

	g.Go(func() error {
		p.Do(ctx)

		return nil
	})

	g.Go(func() error {
		s.Run(ctx)

		return nil
	})

	return g.Wait()
}

func parseLogLevel(s string) *slog.LevelVar {
	var logLevel slog.LevelVar

	switch strings.ToUpper(s) {
	case slog.LevelDebug.String():
		logLevel.Set(slog.LevelDebug)
	case slog.LevelInfo.String():
		logLevel.Set(slog.LevelInfo)
	case slog.LevelWarn.String():
		logLevel.Set(slog.LevelWarn)
	case slog.LevelError.String():
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}

	return &logLevel
}
