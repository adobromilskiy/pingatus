package pinger

import (
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/go-ping/ping"
)

func TestICMPPing(t *testing.T) {
	cfg := &config.EndpointConfig{
		Name:        "test",
		Type:        "icmp",
		Address:     "8.8.8.8",
		PacketCount: 1,
		Interval:    1 * time.Second,
	}

	pinger, _ := ping.NewPinger(cfg.Address)

	p, err := NewICMPPinger(cfg, pinger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	endpoint, err := p.Ping()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Name != cfg.Name {
		t.Errorf("Expected %v, got %v", cfg.Name, endpoint.Name)
	}

	if endpoint.URL != cfg.Address {
		t.Errorf("Expected %v, got %v", cfg.Address, endpoint.URL)
	}

	if endpoint.Status != true {
		t.Errorf("Expected %v, got %v", true, endpoint.Status)
	}
}

// func TestICMPPingError(t *testing.T) {
// 	// Test for a bad host
// 	cfg := &config.EndpointConfig{
// 		Name:        "test",
// 		Type:        "icmp",
// 		Address:     "",
// 		PacketCount: 3,
// 		Interval:    1 * time.Second,
// 	}

// 	_, err := ping.NewPinger(cfg.Address)
// 	if err == nil {
// 		t.Fatalf("Expected error, got %v", err)
// 	}
// }
