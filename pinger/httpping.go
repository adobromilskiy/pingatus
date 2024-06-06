package pinger

import (
	"context"
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/config"
)

type HttpPinger struct {
	Cfg *config.HTTPpointConfig
}

func NewHttpPinger(cfg *config.HTTPpointConfig) *HttpPinger {
	return &HttpPinger{cfg}
}

func (p *HttpPinger) Do(ctx context.Context) {
	ticker := time.NewTicker(p.Cfg.Interval)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("HttpPinger: context is done")
			return
		case <-ticker.C:
			fmt.Println("HttpPinger: tick", p.Cfg.Name)
		}
	}
}
