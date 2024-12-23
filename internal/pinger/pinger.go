package pinger

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/adobromilskiy/pingatus/core"
	"github.com/adobromilskiy/pingatus/internal/config"
	"github.com/adobromilskiy/pingatus/internal/notifier"
)

type pinger interface {
	ping(ctx context.Context) (core.Endpoint, error)
}

type Pingatus struct {
	lg       *slog.Logger
	cfg      []config.EndpointConfig
	storage  core.Setter
	notifier notifier.Notifier
}

func NewPingatus(lg *slog.Logger, cfg []config.EndpointConfig, s core.Setter, n notifier.Notifier) *Pingatus {
	return &Pingatus{lg.With("pkg", "pinger"), cfg, s, n}
}

func (p *Pingatus) Do(ctx context.Context) {
	var wg sync.WaitGroup
	for _, e := range p.cfg {
		wg.Add(1)

		go func(cfg config.EndpointConfig) {
			defer wg.Done()

			switch cfg.Type {
			case "http":
				pinger, err := newHTTP(cfg)
				if err != nil {
					p.lg.Log(
						ctx, slog.LevelError, "can not create pinger",
						"type", cfg.Type,
						"endpoint", cfg.Name,
						"err", err,
					)

					return
				}

				p.lg.Log(
					ctx, slog.LevelInfo, "start pinger",
					"type", cfg.Type,
					"endpoint", cfg.Name,
				)

				p.run(ctx, pinger, cfg)
			case "icmp":
				pinger, err := newICMP(cfg)
				if err != nil {
					p.lg.Log(
						ctx, slog.LevelError, "can not create pinger",
						"type", cfg.Type,
						"endpoint", cfg.Name,
						"err", err,
					)

					return
				}

				p.lg.Log(
					ctx, slog.LevelInfo, "start pinger",
					"type", cfg.Type,
					"endpoint", cfg.Name,
				)

				p.run(ctx, pinger, cfg)
			}
		}(e)
	}

	wg.Wait()

	p.lg.Log(ctx, slog.LevelInfo, "all pingers stopped")
}

func (p *Pingatus) run(ctx context.Context, pinger pinger, cfg config.EndpointConfig) {
	ticker := time.NewTicker(cfg.Interval)

	currentStatus := true

	for {
		select {
		case <-ctx.Done():
			p.lg.Log(ctx, slog.LevelInfo, "pinging stopped via context", "endpoint", cfg.Name)

			return
		case <-ticker.C:
			p.lg.Log(ctx, slog.LevelDebug, "pinging", "endpoint", cfg.Name)

			endpoint, err := pinger.ping(ctx)
			if err != nil {
				p.lg.Log(
					ctx, slog.LevelWarn, "got error while pinging",
					"type", cfg.Type,
					"endpoint", cfg.Name,
					"err", err,
				)

				continue
			}

			if endpoint.Status && !currentStatus {
				currentStatus = true

				go p.notifier.Send(ctx, "endpoint "+cfg.Name+" is online")
			}

			if !endpoint.Status && currentStatus {
				currentStatus = false

				go p.notifier.Send(ctx, "endpoint "+cfg.Name+" is offline")
			}

			err = p.storage.Save(endpoint)
			if err != nil {
				p.lg.Log(
					ctx, slog.LevelWarn, "got error while saving endpoint",
					"type", cfg.Type,
					"endpoint", cfg.Name,
					"err", err,
				)
			}
		}
	}
}
