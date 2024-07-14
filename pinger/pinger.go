package pinger

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/storage"
	"github.com/go-ping/ping"
)

type Pinger interface {
	Ping() (*storage.Endpoint, error)
}

type Pingatus struct {
	Cfg           *config.Config
	CurrentStatus bool
	Storage       storage.Storage
	Notifier      notifier.Notifier
}

func NewPingatus(cfg *config.Config, s storage.Storage, n notifier.Notifier) *Pingatus {
	return &Pingatus{cfg, true, s, n}
}

func (p *Pingatus) Do(ctx context.Context) {

	var wg sync.WaitGroup
	for _, e := range p.Cfg.EndPoints {
		wg.Add(1)
		go func(cfg *config.EndpointConfig) {
			defer wg.Done()
			switch cfg.Type {
			case "http":
				pinger, err := NewHTTPPinger(cfg)
				if err != nil {
					log.Printf("[ERROR] HTTPPinger %s: error creating pinger: %v", cfg.Name, err)
					return
				}
				p.Run(ctx, pinger, cfg)
			case "icmp":
				icmppinger, err := ping.NewPinger(cfg.Address)
				if err != nil {
					log.Printf("[ERROR] ICMPPinger %s: error creating pinger: %v", cfg.Name, err)
					return
				}

				pinger, err := NewICMPPinger(cfg, icmppinger)
				if err != nil {
					log.Printf("[ERROR] ICMPPinger %s: error creating pinger: %v", cfg.Name, err)
					return
				}

				p.Run(ctx, pinger, cfg)
			}
		}(&e)
	}
	wg.Wait()
	log.Println("[INFO] all pings finished")
}

func (p *Pingatus) Run(ctx context.Context, pinger Pinger, cfg *config.EndpointConfig) {
	ticker := time.NewTicker(cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] %s pinger '%s': stoped via context", cfg.Type, cfg.Name)
			return
		case <-ticker.C:
			endpoint, err := pinger.Ping()
			if err != nil {
				log.Printf("[ERROR] %s pinger '%s': error pinging: %v", cfg.Type, cfg.Name, err)
				continue
			}
			if endpoint.Status && !p.CurrentStatus {
				p.CurrentStatus = true
				go p.Notifier.Send("endpoint " + cfg.Name + " is online")
			}
			if !endpoint.Status && p.CurrentStatus {
				p.CurrentStatus = false
				go p.Notifier.Send("endpoint " + cfg.Name + " is offline")
			}
			err = p.Storage.SaveEndpoint(ctx, endpoint)
			if err != nil {
				log.Printf("[ERROR] %s pinger '%s': error save endpoint: %v", cfg.Type, cfg.Name, err)
			}
		}
	}
}
