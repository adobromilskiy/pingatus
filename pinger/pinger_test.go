package pinger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/mock"
)

func TestPingerDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		HTTPPoint: []config.HTTPpointConfig{{
			Name:     "test",
			URL:      server.URL,
			Status:   http.StatusOK,
			Interval: 300 * time.Millisecond,
			Timeout:  100 * time.Millisecond,
		}},
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewPingatus(cfg, store, notifier)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	pinger.Do(ctx)
}
