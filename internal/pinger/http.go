package pinger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/core"
	"github.com/adobromilskiy/pingatus/internal/config"
)

type httpPinger struct {
	cfg config.EndpointConfig
}

func newHTTP(cfg config.EndpointConfig) (*httpPinger, error) {
	if cfg.Status == 0 {
		return nil, errStatusNotSet
	}

	if cfg.Timeout == 0 {
		return nil, errTimeoutNotSet
	}

	return &httpPinger{cfg}, nil
}

func (p *httpPinger) ping(ctx context.Context) (core.Endpoint, error) {
	client := &http.Client{
		Timeout: p.cfg.Timeout,
	}

	endpoint := core.Endpoint{
		Name:    p.cfg.Name,
		Address: p.cfg.Address,
		Date:    time.Now().Unix(),
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.Address, nil)
	if err != nil {
		return endpoint, fmt.Errorf("%w: %w", errFailedCreateRequest, err)
	}

	resp, err := client.Do(r)
	if err != nil {
		endpoint.Status = false

		return endpoint, nil
	}
	defer resp.Body.Close()

	endpoint.Status = resp.StatusCode == p.cfg.Status

	return endpoint, nil
}
