package pinger

import (
	"context"
	"log"
	"net/http"
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
			log.Printf("[INFO] HttpPinger %s: stoped via context", p.Cfg.Name)
			return
		case <-ticker.C:
			p.ping()
		}
	}
}

func (p *HttpPinger) ping() {
	client := &http.Client{
		Timeout: p.Cfg.Timeout,
	}
	req, err := http.NewRequest("GET", p.Cfg.URL, nil)
	if err != nil {
		log.Printf("[ERROR] HttpPinger %s: error creating request: %v", p.Cfg.Name, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] HttpPinger %s: error sending request: %v", p.Cfg.Name, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("HttpPinger %s: received response: %d\n", p.Cfg.Name, resp.StatusCode)
}
