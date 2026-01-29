package core

import "context"

// Endpoint represents a single health check result.
type Endpoint struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
	Date    int64  `json:"date"`
}

// Setter defines storage operations for endpoint checks.
type Setter interface {
	Save(ctx context.Context, e Endpoint) error
	GetEndpoints(ctx context.Context) ([]string, error)
	GetEndpointStats(ctx context.Context, name string, date int64) ([]Endpoint, error)
	GetLastSuccess(ctx context.Context, name string) (*Endpoint, error)
	GetLastFailure(ctx context.Context, name string) (*Endpoint, error)
}
