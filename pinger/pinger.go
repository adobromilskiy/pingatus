package pinger

import (
	"context"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/storage"
)

type Pinger struct {
	Cfg      *config.Config
	Storage  storage.Storage
	Notifier notifier.Notifier
}

func NewPinger(cfg *config.Config, s storage.Storage, n notifier.Notifier) *Pinger {
	return &Pinger{cfg, s, n}
}

func (p *Pinger) Do(ctx context.Context) {

	var wg sync.WaitGroup
	for _, httpPoint := range p.Cfg.HTTPPoint {
		wg.Add(1)
		go func(p *config.HTTPpointConfig, s storage.Storage, n notifier.Notifier) {
			defer wg.Done()
			NewHTTPPinger(p, s, n).Do(ctx)
		}(&httpPoint, p.Storage, p.Notifier)
	}
	wg.Wait()
	log.Println("[INFO] all pings finished")
}
