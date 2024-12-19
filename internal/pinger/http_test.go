package pinger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/internal/config"
)

func TestHTTPPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.EndpointConfig{
		Name:     "test",
		Address:  server.URL,
		Status:   http.StatusOK,
		Interval: 1 * time.Second,
		Timeout:  1 * time.Second,
	}

	pinger, err := newHTTP(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ctx := context.Background()

	endpoint, err := pinger.ping(ctx)
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

	// Test for a non-200 status code
	server.Config.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
	})

	endpoint, err = pinger.ping(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Status != false {
		t.Errorf("Expected %v, got %v", false, endpoint.Status)
	}

	// Test for a request error
	server.Close()

	_, err = pinger.ping(ctx)
	if err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestHTTPPing_Error(t *testing.T) {
	cfg := config.EndpointConfig{
		Name:     "test",
		Address:  "adsfsdf://invalid_url()!$@*(^*)",
		Status:   http.StatusOK,
		Interval: 1 * time.Second,
		Timeout:  1 * time.Second,
	}

	pinger, err := newHTTP(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = pinger.ping(context.Background())
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestHTTPPing_ErrorCFG(t *testing.T) {
	cfg := config.EndpointConfig{
		Name:     "test",
		Type:     "http",
		Address:  "http://example.com",
		Status:   http.StatusOK,
		Interval: 200 * time.Millisecond,
		Timeout:  0,
	}

	_, err := newHTTP(cfg)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}

	cfg = config.EndpointConfig{
		Name:     "test",
		Type:     "http",
		Address:  "http://example.com",
		Status:   0,
		Interval: 200 * time.Millisecond,
		Timeout:  0,
	}

	_, err = newHTTP(cfg)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}
