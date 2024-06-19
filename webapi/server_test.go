package webapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/mock"
)

func TestServer_Run(t *testing.T) {
	cfg := config.WEBAPIConfig{
		ListenAddr: ":8080",
		AssetsDir:  "./",
	}
	store := &mock.StoreMock{}
	server := NewServer(cfg, store)

	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx)

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request to the server
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", resp.StatusCode)
	}

	// Stop the server
	cancel()

	// Give the server some time to stop
	time.Sleep(1 * time.Second)

	// Make a request to the server, expect an error because the server should be stopped
	_, err = http.Get("http://localhost:8080")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
