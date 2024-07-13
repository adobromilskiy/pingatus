package pinger

import (
	"context"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/mock"
	"github.com/go-ping/ping"
)

func TestICMPPing(t *testing.T) {
	cfg := &config.ICMPpointConfig{
		Name:        "test",
		IP:          "8.8.8.8",
		PacketCount: 3,
		Interval:    1 * time.Second,
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}
	pinger, _ := ping.NewPinger(cfg.IP)

	p := NewICMPPinger(cfg, store, notifier, pinger)

	endpoint, err := p.ping()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Name != cfg.Name {
		t.Errorf("Expected %v, got %v", cfg.Name, endpoint.Name)
	}

	if endpoint.URL != cfg.IP {
		t.Errorf("Expected %v, got %v", cfg.IP, endpoint.URL)
	}

	if endpoint.Status != true {
		t.Errorf("Expected %v, got %v", true, endpoint.Status)
	}
}

func TestICMPPingError(t *testing.T) {
	// Test for a bad host
	cfg := &config.ICMPpointConfig{
		Name:        "test",
		IP:          "",
		PacketCount: 3,
		Interval:    1 * time.Second,
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}
	pinger, _ := ping.NewPinger(cfg.IP)

	p := NewICMPPinger(cfg, store, notifier, pinger)

	_, err := p.ping()
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func TestICMPDo(t *testing.T) {
	cfg := &config.ICMPpointConfig{
		Name:        "test",
		IP:          "8.8.8.8",
		PacketCount: 1,
		Interval:    1 * time.Second,
	}
	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}
	pinger, _ := ping.NewPinger(cfg.IP)
	p := NewICMPPinger(cfg, store, notifier, pinger)

	// Test if Do stops when context is done
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	p.Do(ctx)
}

func TestICMPDoPingError(t *testing.T) {
	// Test for a bad host
	cfg := &config.ICMPpointConfig{
		Name:        "test",
		IP:          "",
		PacketCount: 1,
		Interval:    1 * time.Second,
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}
	pinger, _ := ping.NewPinger(cfg.IP)

	p := NewICMPPinger(cfg, store, notifier, pinger)

	// Test if Do stops when context is done
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	p.Do(ctx)
}

func TestICMPDoChangeStatus(t *testing.T) {
	cfg := &config.ICMPpointConfig{
		Name:        "test",
		IP:          "8.8.8.8",
		PacketCount: 1,
		Interval:    1 * time.Second,
	}
	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}
	pinger, _ := ping.NewPinger(cfg.IP)
	p := NewICMPPinger(cfg, store, notifier, pinger)
	p.CurrentStatus = false

	// Test if Do stops when context is done
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	p.Do(ctx)
}
