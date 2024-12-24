package core

import "context"

type Endpoint struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
	Date    int64  `json:"date"`
}

type Setter interface {
	Save(ctx context.Context, e Endpoint) error
	GetEndpoints(ctx context.Context) ([]string, error)
	GetEndpointStats(ctx context.Context, name string, date int64) ([]Endpoint, error)
}
