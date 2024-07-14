package pinger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/storage"
)

type HTTPPinger struct {
	Cfg *config.EndpointConfig
}

func NewHTTPPinger(cfg *config.EndpointConfig) (*HTTPPinger, error) {
	if cfg.Status == 0 {
		return nil, fmt.Errorf("status is not set")
	}

	if cfg.Timeout == 0 {
		return nil, fmt.Errorf("timeout is not set")
	}

	return &HTTPPinger{cfg}, nil
}

func (p *HTTPPinger) Ping() (*storage.Endpoint, error) {
	client := &http.Client{
		Timeout: p.Cfg.Timeout,
	}
	req, err := http.NewRequest("GET", p.Cfg.Address, nil)
	if err != nil {
		return nil, err
	}

	endpoint := storage.Endpoint{
		Name: p.Cfg.Name,
		URL:  p.Cfg.Address,
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
