package notifier

import (
	"sync"
	"testing"

	"github.com/adobromilskiy/pingatus/config"
)

func TestGet(t *testing.T) {
	cfg := &config.Config{
		Notifier: config.NotifierConfig{
			Type:     "telegram",
			TgToken:  "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			TgChatID: "-1001234567890",
		},
	}

	_, err := Get(cfg)
	if err != nil {
		t.Fatalf("Failed to get notifier: %v", err)
	}
}

func TestGet_ErrorEmptyToken(t *testing.T) {
	notifierOnce = sync.Once{}

	cfg := &config.Config{
		Notifier: config.NotifierConfig{
			Type:     "telegram",
			TgChatID: "-1001234567890",
		},
	}

	_, err := Get(cfg)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGet_ErrorEmptyChatID(t *testing.T) {
	notifierOnce = sync.Once{}

	cfg := &config.Config{
		Notifier: config.NotifierConfig{
			Type:    "telegram",
			TgToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		},
	}

	_, err := Get(cfg)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGet_ErrorUnknownType(t *testing.T) {
	notifierOnce = sync.Once{}

	cfg := &config.Config{
		Notifier: config.NotifierConfig{
			Type: "unknown",
		},
	}

	_, err := Get(cfg)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
