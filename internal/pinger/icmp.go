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
	cfg    config.EndpointConfig
	pinger *probing.Pinger
}

func newICMP(cfg config.EndpointConfig) (*icmpPinger, error) {
	if cfg.PacketCount == 0 {
		return nil, errPacketCountNotSet
	}

	if cfg.Timeout == 0 {
		return nil, errTimeoutNotSet
	}

	pinger, err := probing.NewPinger(cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("pinger: %w", err)
	}

	return &icmpPinger{
		cfg:    cfg,
		pinger: pinger,
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

	go func() {
		<-ctx.Done()
		p.pinger.Stop()
	}()

	p.pinger.Count = p.cfg.PacketCount

	if err := p.pinger.Run(); err != nil {
		return endpoint, err
	}

	stats := p.pinger.Statistics()

	endpoint.Status = stats.PacketsRecv > 0

	return endpoint, nil
}
