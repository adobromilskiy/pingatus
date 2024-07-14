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
		EndPoints: []config.EndpointConfig{{
			Name:     "test",
			Type:     "http",
			Address:  server.URL,
			Status:   http.StatusOK,
			Interval: 200 * time.Millisecond,
			Timeout:  100 * time.Millisecond,
		}, {
			Name:        "test2",
			Type:        "icmp",
			Address:     "8.8.8.8",
			PacketCount: 1,
			Interval:    100 * time.Millisecond,
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

func TestPingerDo_errors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		EndPoints: []config.EndpointConfig{{
			Name:     "test",
			Type:     "http",
			Address:  server.URL,
			Interval: 20 * time.Millisecond,
			Timeout:  10 * time.Millisecond,
		}, {
			Name:        "test2",
			Type:        "icmp",
			Address:     "8.8.8.8",
			PacketCount: 0,
			Interval:    10 * time.Millisecond,
		}},
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewPingatus(cfg, store, notifier)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	pinger.Do(ctx)
}

func TestPingerDo_invalidIP(t *testing.T) {
	cfg := &config.Config{
		EndPoints: []config.EndpointConfig{{
			Name:        "test2",
			Type:        "icmp",
			Address:     "",
			PacketCount: 3,
			Interval:    10 * time.Millisecond,
		}},
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewPingatus(cfg, store, notifier)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	pinger.Do(ctx)
}

func TestPingDo_changeStatus(t *testing.T) {
	cfg := &config.Config{
		EndPoints: []config.EndpointConfig{{
			Name:        "test2",
			Type:        "icmp",
			Address:     "8.8.8.8",
			PacketCount: 1,
			Interval:    100 * time.Millisecond,
		}},
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewPingatus(cfg, store, notifier)
	pinger.CurrentStatus = false

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	pinger.Do(ctx)
}
