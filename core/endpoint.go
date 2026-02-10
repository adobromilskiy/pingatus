package core

import "context"

// Endpoint represents a single availability datapoint.
type Endpoint struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
	Date    int64  `json:"date"`
}

// EndpointReader reads endpoint names and stats.
type EndpointReader interface {
	GetEndpoints(ctx context.Context) ([]string, error)
	GetEndpointStats(ctx context.Context, name string, from int64, to int64) ([]Endpoint, error)
}

// EndpointWriter persists endpoint stats.
type EndpointWriter interface {
	Save(ctx context.Context, e Endpoint) error
}

// EndpointStore provides read/write access to endpoint stats.
type EndpointStore interface {
	EndpointReader
	EndpointWriter
}
