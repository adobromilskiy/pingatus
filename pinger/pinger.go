package pinger

import (
	"context"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
)

type Pinger struct {
	Cfg *config.Config
}

func NewPinger(cfg *config.Config) *Pinger {
	return &Pinger{cfg}
}

func (p *Pinger) Do(ctx context.Context) {

	var wg sync.WaitGroup
	for _, httpPoint := range p.Cfg.HTTPPoint {
		wg.Add(1)
		go func(p *config.HTTPpointConfig) {
			defer wg.Done()
			NewHTTPPinger(p).Do(ctx)
		}(&httpPoint)
	}
	wg.Wait()
	log.Println("[INFO] all pings finished")
}
