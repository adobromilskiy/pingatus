package pinger

import (
	"context"

	"github.com/adobromilskiy/pingatus/core"
)

type Pinger interface {
	ping(ctx context.Context) (core.Endpoint, error)
}
