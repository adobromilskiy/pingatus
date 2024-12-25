package pinger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/adobromilskiy/pingatus/core"
	"github.com/adobromilskiy/pingatus/internal/config"
	"github.com/adobromilskiy/pingatus/internal/notifier"
)

var (
	errStatusNotSet        = errors.New("pinger: status is not set")
	errTimeoutNotSet       = errors.New("pinger: timeout is not set")
	errFailedCreateRequest = errors.New("pinger: failed to create request")
	errPacketCountNotSet   = errors.New("pinger: packetcount is not set")
)

type pinger interface {
	ping(ctx context.Context) (core.Endpoint, error)
}

type Pingatus struct {
	lg       *slog.Logger
	cfg      []config.EndpointConfig
	storage  core.Setter
	notifier notifier.Notifier
	mu       sync.Mutex
	status   map[string]bool
}

func NewPingatus(lg *slog.Logger, cfg []config.EndpointConfig, s core.Setter, n notifier.Notifier) *Pingatus {
	return &Pingatus{
		lg:       lg.With("pkg", "pinger"),
		cfg:      cfg,
		storage:  s,
		notifier: n,
		status:   make(map[string]bool, len(cfg)),
	}
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

				p.mu.Lock()
				p.status[cfg.Name] = true
				p.mu.Unlock()

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

				p.mu.Lock()
				p.status[cfg.Name] = true
				p.mu.Unlock()

				p.run(ctx, pinger, cfg)
			}
		}(e)
	}

	wg.Wait()

	p.lg.Log(ctx, slog.LevelInfo, "all pingers stopped")
}

func (p *Pingatus) run(ctx context.Context, pinger pinger, cfg config.EndpointConfig) {
	ticker := time.NewTicker(cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			p.lg.Log(ctx, slog.LevelInfo, "pinging stopped via context", "endpoint", cfg.Name)

			return
		case <-ticker.C:
			p.mu.Lock()
			p.lg.Log(ctx, slog.LevelDebug, "pinging", "endpoint", cfg.Name, "status", p.status[cfg.Name])
			p.mu.Unlock()

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

			fmt.Println(endpoint)

			p.mu.Lock()
			if endpoint.Status && !p.status[cfg.Name] {
				p.status[cfg.Name] = true

				go p.notifier.Send(ctx, "endpoint "+cfg.Name+" is online")
			}

			if !endpoint.Status && p.status[cfg.Name] {
				p.status[cfg.Name] = false

				go p.notifier.Send(ctx, "endpoint "+cfg.Name+" is offline")
			}
			p.mu.Unlock()

			err = p.storage.Save(ctx, endpoint)
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
