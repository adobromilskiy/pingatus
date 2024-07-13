package pinger

import (
	"context"
	"log"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/storage"
	"github.com/go-ping/ping"
)

type ICMPPinger struct {
	Cfg           *config.ICMPpointConfig
	CurrentStatus bool
	Storage       storage.Storage
	Notifier      notifier.Notifier
	Pinger        *ping.Pinger
}

func NewICMPPinger(cfg *config.ICMPpointConfig, s storage.Storage, n notifier.Notifier, p *ping.Pinger) *ICMPPinger {
	return &ICMPPinger{
		Cfg:           cfg,
		CurrentStatus: true,
		Storage:       s,
		Notifier:      n,
		Pinger:        p,
	}
}

func (p *ICMPPinger) Do(ctx context.Context) {
	ticker := time.NewTicker(p.Cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] ICMPPinger %s: stoped via context", p.Cfg.Name)
			return
		case <-ticker.C:
			endpoint, err := p.ping()
			if err != nil {
				log.Printf("[ERROR] ICMPPinger %s: error pinging: %v", p.Cfg.Name, err)
				continue
			}
			if endpoint.Status && !p.CurrentStatus {
				p.CurrentStatus = true
				go p.Notifier.Send("endpoint " + p.Cfg.Name + " is online")
			}
			if !endpoint.Status && p.CurrentStatus {
				p.CurrentStatus = false
				go p.Notifier.Send("endpoint " + p.Cfg.Name + " is offline")
			}
			err = p.Storage.SaveEndpoint(ctx, endpoint)
			if err != nil {
				log.Printf("[ERROR] ICMPPinger %s: error save endpoint: %v", p.Cfg.Name, err)
			}
		}
	}
}

func (p *ICMPPinger) ping() (*storage.Endpoint, error) {
	p.Pinger.Count = p.Cfg.PacketCount
	if err := p.Pinger.Run(); err != nil {
		return nil, err

	}

	stats := p.Pinger.Statistics()

	endpoint := storage.Endpoint{
		Name: p.Cfg.Name,
		URL:  p.Cfg.IP,
		Date: time.Now().Unix(),
	}

	endpoint.Status = stats.PacketsRecv > 0

	return &endpoint, nil
}
