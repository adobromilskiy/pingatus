package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetLastEndpoint(ctx context.Context, filter primitive.M) (*Endpoint, error)
	GetEndpoints(ctx context.Context, filter primitive.M) ([]*Endpoint, error)
	SaveEndpoint(ctx context.Context, endpoint *Endpoint) error
	Close()
}
