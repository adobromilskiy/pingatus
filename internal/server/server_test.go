package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/internal/mocks"
)

func TestServer_Run(t *testing.T) {
	logger := mocks.MockLogger()
	store := &mocks.StorageMock{}
	server := New(logger, store, ":8081")

	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx)

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request to the server
	resp, err := http.Get("http://localhost:8081/endpoints")
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
	_, err = http.Get("http://localhost:8081")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
