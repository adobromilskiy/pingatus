package pinger

import "github.com/adobromilskiy/pingatus/core"

type Pinger interface {
	Ping() (core.Endpoint, error)
}
