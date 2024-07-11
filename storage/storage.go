package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetLastEndpoint(context.Context, primitive.M) (*Endpoint, error)
	GetEndpoints(context.Context, primitive.M) ([]*Endpoint, error)
	SaveEndpoint(context.Context, *Endpoint) error
	GetNames(context.Context) ([]string, error)
	Close(context.Context)
}
