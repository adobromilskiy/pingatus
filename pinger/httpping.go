package pinger

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/notifier"
	"github.com/adobromilskiy/pingatus/storage"
)

type HTTPPinger struct {
	Cfg           *config.HTTPpointConfig
	CurrentStatus bool
	Storage       storage.Storage
	Notifier      notifier.Notifier
}

func NewHTTPPinger(cfg *config.HTTPpointConfig, s storage.Storage, n notifier.Notifier) *HTTPPinger {
	return &HTTPPinger{cfg, true, s, n}
}

func (p *HTTPPinger) Do(ctx context.Context) {
	ticker := time.NewTicker(p.Cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] HTTPPinger %s: stoped via context", p.Cfg.Name)
			return
		case <-ticker.C:
			endpoint, err := p.ping()
			if err != nil {
				log.Printf("[ERROR] HTTPPinger %s: error pinging: %v", p.Cfg.Name, err)
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
				log.Printf("[ERROR] HTTPPinger %s: error save endpoint: %v", p.Cfg.Name, err)
			}
		}
	}
}

func (p *HTTPPinger) ping() (*storage.Endpoint, error) {
	client := &http.Client{
		Timeout: p.Cfg.Timeout,
	}
	req, err := http.NewRequest("GET", p.Cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	endpoint := storage.Endpoint{
		Name: p.Cfg.Name,
		URL:  p.Cfg.URL,
		Date: time.Now().Unix(),
	}

	resp, err := client.Do(req)
	if err != nil {
		endpoint.Status = false
		return &endpoint, nil
	}
	defer resp.Body.Close()

	endpoint.Status = resp.StatusCode == p.Cfg.Status

	return &endpoint, nil
}
