package pinger

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/storage"
)

type HTTPPinger struct {
	Cfg *config.HTTPpointConfig
}

func NewHTTPPinger(cfg *config.HTTPpointConfig) *HTTPPinger {
	return &HTTPPinger{cfg}
}

func (p *HTTPPinger) Do(ctx context.Context) {
	ticker := time.NewTicker(p.Cfg.Interval)
	store, err := storage.GetMongoClient()
	if err != nil {
		log.Printf("[ERROR] HTTPPinger %s: error getting mongo client: %v", p.Cfg.Name, err)
		return
	}
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
			err = store.SaveEndpoint(ctx, endpoint)
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
