package notifier

import (
	"testing"

	"github.com/adobromilskiy/pingatus/internal/config"
	"github.com/adobromilskiy/pingatus/internal/mocks"
)

func TestNew_Success(t *testing.T) {
	logger := mocks.MockLogger()

	cfg := config.NotifierConfig{
		Type:     "telegram",
		TgToken:  "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		TgChatID: "-1001234567890",
	}

	_, err := New(logger, cfg)
	if err != nil {
		t.Fatalf("Failed to get notifier: %v", err)
	}
}

func TestNew_ErrorEmptyToken(t *testing.T) {
	logger := mocks.MockLogger()

	cfg := config.NotifierConfig{
		Type:     "telegram",
		TgChatID: "-1001234567890",
	}

	_, err := New(logger, cfg)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestNew_ErrorEmptyChatID(t *testing.T) {
	logger := mocks.MockLogger()

	cfg := config.NotifierConfig{
		Type:    "telegram",
		TgToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
	}

	_, err := New(logger, cfg)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestNew_ErrorUnknownType(t *testing.T) {
	logger := mocks.MockLogger()

	cfg := config.NotifierConfig{
		Type: "unknown",
	}

	_, err := New(logger, cfg)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
