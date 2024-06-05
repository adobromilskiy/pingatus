package pinger

import (
	"context"
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/config"
)

type HttpPinger struct {
	Url        string
	HttpStatus int
	Timeout    time.Duration
	Interval   time.Duration
}

func NewHttpPinger(url string, httpStatus int, timeout, interval time.Duration) *HttpPinger {
	return &HttpPinger{
		Url:        url,
		HttpStatus: httpStatus,
		Timeout:    timeout,
		Interval:   interval,
	}
}

func (p *HttpPinger) Do(ctx context.Context) {
	if _, err := config.Load(); err != nil {
		fmt.Println("Failed to load config:", err)
	}
	ticker := time.NewTicker(p.Interval)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("HttpPinger: context is done")
			return
		case <-ticker.C:
			fmt.Println("HttpPinger: tick")
		}
	}
}
