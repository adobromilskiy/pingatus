package pinger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/core"
	"github.com/adobromilskiy/pingatus/internal/config"
	probing "github.com/prometheus-community/pro-bing"
)

var errPacketCountNotSet = errors.New("pinger: packetcount is not set")

type icmpPinger struct {
	cfg    config.EndpointConfig
	pinger *probing.Pinger
}

func newICMP(cfg config.EndpointConfig) (*icmpPinger, error) {
	if cfg.PacketCount == 0 {
		return nil, errPacketCountNotSet
	}

	pinger, err := probing.NewPinger(cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("pinger: %w", err)
	}

	pinger.Count = cfg.PacketCount

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

	if err := p.pinger.Run(); err != nil {
		return endpoint, err
	}

	stats := p.pinger.Statistics()

	endpoint.Status = stats.PacketsRecv > 0

	return endpoint, nil
}