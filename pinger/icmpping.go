package pinger

import (
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/storage"
	"github.com/go-ping/ping"
)

type ICMPPinger struct {
	Cfg    *config.EndpointConfig
	Pinger *ping.Pinger
}

func NewICMPPinger(cfg *config.EndpointConfig, p *ping.Pinger) (*ICMPPinger, error) {
	if cfg.PacketCount == 0 {
		return nil, fmt.Errorf("packetcount is not set")
	}

	return &ICMPPinger{
		Cfg:    cfg,
		Pinger: p,
	}, nil
}

func (p *ICMPPinger) Ping() (*storage.Endpoint, error) {
	p.Pinger.Count = p.Cfg.PacketCount
	if err := p.Pinger.Run(); err != nil {
		return nil, err

	}

	stats := p.Pinger.Statistics()

	endpoint := storage.Endpoint{
		Name: p.Cfg.Name,
		URL:  p.Cfg.Address,
		Date: time.Now().Unix(),
	}

	endpoint.Status = stats.PacketsRecv > 0

	return &endpoint, nil
}
