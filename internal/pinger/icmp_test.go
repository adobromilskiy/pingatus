package pinger

import (
	"context"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/internal/config"
)

func TestICMPPing(t *testing.T) {
	cfg := config.EndpointConfig{
		Name:        "test",
		Type:        "icmp",
		Address:     "8.8.8.8",
		PacketCount: 1,
		Timeout:     time.Second,
		Interval:    1 * time.Second,
	}

	p, err := newICMP(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	endpoint, err := p.ping(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Name != cfg.Name {
		t.Errorf("Expected %v, got %v", cfg.Name, endpoint.Name)
	}

	if endpoint.Address != cfg.Address {
		t.Errorf("Expected %v, got %v", cfg.Address, endpoint.Address)
	}

	if endpoint.Status != true {
		t.Errorf("Expected %v, got %v", true, endpoint.Status)
	}
}

func TestICMPPing_Error(t *testing.T) {
	cfg := config.EndpointConfig{
		Name:        "test",
		Type:        "icmp",
		Address:     "",
		PacketCount: 3,
		Timeout:     time.Second,
		Interval:    1 * time.Second,
	}

	p, err := newICMP(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = p.ping(context.Background())
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}

	cfg = config.EndpointConfig{
		Name:        "test",
		Type:        "icmp",
		Address:     "",
		PacketCount: 0,
		Timeout:     time.Second,
		Interval:    1 * time.Second,
	}

	_, err = newICMP(cfg)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}

	cfg = config.EndpointConfig{
		Name:        "test",
		Type:        "icmp",
		Address:     "",
		PacketCount: 3,
		Interval:    1 * time.Second,
	}

	_, err = newICMP(cfg)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}
