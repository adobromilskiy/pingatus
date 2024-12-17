package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/adobromilskiy/pingatus/internal/config"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	logLevel := parseLogLevel(cfg.Logger.Level)

	lg := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	if cfg.Logger.IsJSON {
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	lg.Log(ctx, slog.LevelInfo, "application started", "version", "2.0")

	return nil
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
