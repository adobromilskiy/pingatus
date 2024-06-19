package pinger

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/mock"
)

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.HTTPpointConfig{
		Name:     "test",
		URL:      server.URL,
		Status:   http.StatusOK,
		Interval: 1 * time.Second,
		Timeout:  1 * time.Second,
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewHTTPPinger(cfg, store, notifier)

	endpoint, err := pinger.ping()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Name != cfg.Name {
		t.Errorf("Expected %v, got %v", cfg.Name, endpoint.Name)
	}

	if endpoint.URL != cfg.URL {
		t.Errorf("Expected %v, got %v", cfg.URL, endpoint.URL)
	}

	if endpoint.Status != true {
		t.Errorf("Expected %v, got %v", true, endpoint.Status)
	}

	// Test for a non-200 status code
	server.Config.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
	})

	endpoint, err = pinger.ping()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpoint.Status != false {
		t.Errorf("Expected %v, got %v", false, endpoint.Status)
	}

	// Test for a request error
	server.Close()

	_, err = pinger.ping()
	if err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestPingError(t *testing.T) {
	cfg := &config.HTTPpointConfig{
		Name:     "test",
		URL:      "adsfsdf://invalid_url()!$@*(^*)",
		Status:   http.StatusOK,
		Interval: 1 * time.Second,
		Timeout:  1 * time.Second,
	}

	store := &mock.StoreMock{}
	notifier := &mock.NotifierMock{}

	pinger := NewHTTPPinger(cfg, store, notifier)

	_, err := pinger.ping()
	if err == nil {
		t.Fatalf("Expected error")
	}
}
