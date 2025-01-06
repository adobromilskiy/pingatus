package pinger

import (
	"context"
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/core"
	"github.com/adobromilskiy/pingatus/internal/config"
	probing "github.com/prometheus-community/pro-bing"
)

type icmpPinger struct {
	cfg config.EndpointConfig
}

func newICMP(cfg config.EndpointConfig) (*icmpPinger, error) {
	if cfg.PacketCount == 0 {
		return nil, errPacketCountNotSet
	}

	if cfg.Timeout == 0 {
		return nil, errTimeoutNotSet
	}

	return &icmpPinger{
		cfg: cfg,
	}, nil
}

func (p *icmpPinger) ping(ctx context.Context) (core.Endpoint, error) {
	endpoint := core.Endpoint{
		Name:    p.cfg.Name,
		Address: p.cfg.Address,
		Date:    time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
	defer cancel()

	pinger, err := probing.NewPinger(p.cfg.Address)
	if err != nil {
		return endpoint, fmt.Errorf("pinger: %w", err)
	}

	pinger.Count = p.cfg.PacketCount

	go func() {
		<-ctx.Done()
		pinger.Stop()
	}()

	if err := pinger.Run(); err != nil {
		return endpoint, err
	}

	stats := pinger.Statistics()

	endpoint.Status = stats.PacketsRecv > 0

	return endpoint, nil
}
