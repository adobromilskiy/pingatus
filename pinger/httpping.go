package pinger

import (
	"context"
	"fmt"
	"time"
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
