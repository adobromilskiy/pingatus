package pinger

import (
	"context"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/storage"
	"github.com/go-ping/ping"
)

type Pinger interface {
	Ping() (*storage.Endpoint, error)
}

type Pingatus struct {
	Cfg      *config.Config
	Storage  storage.Storage
	Notifier notifier.Notifier
}

func NewPingatus(cfg *config.Config, s storage.Storage, n notifier.Notifier) *Pingatus {
	return &Pingatus{cfg, s, n}
}

func (p *Pingatus) Do(ctx context.Context) {

	var wg sync.WaitGroup
	for _, httpPoint := range p.Cfg.HTTPPoint {
		wg.Add(1)
		go func(p *config.HTTPpointConfig, s storage.Storage, n notifier.Notifier) {
			defer wg.Done()
			NewHTTPPinger(p, s, n).Do(ctx)
		}(&httpPoint, p.Storage, p.Notifier)
	}
	for _, icmpPoint := range p.Cfg.ICMPPoint {
		wg.Add(1)
		go func(p *config.ICMPpointConfig, s storage.Storage, n notifier.Notifier) {
			defer wg.Done()
			pinger, err := ping.NewPinger(p.IP)
			if err != nil {
				log.Printf("[ERROR] ICMPPinger %s: error creating pinger: %v", p.Name, err)
				return
			}
			NewICMPPinger(p, s, n, pinger).Do(ctx)
		}(&icmpPoint, p.Storage, p.Notifier)
	}
	wg.Wait()
	log.Println("[INFO] all pings finished")
}
