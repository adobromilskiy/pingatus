package notifier

import (
	"os"
	"sync"
	"testing"

	"github.com/adobromilskiy/pingatus/config"
	"gopkg.in/yaml.v3"
)

func TestGet(t *testing.T) {
	defer config.Reset()
	expectedConfig := &config.Config{
		Notifier: config.NotifierConfig{
			Type:     "telegram",
			TgToken:  "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			TgChatID: "-1001234567890",
		},
	}

	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

	data, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = Get()
	if err != nil {
		t.Fatalf("Failed to get notifier: %v", err)
	}
}

func TestGet_ErrorEmptyToken(t *testing.T) {
	defer config.Reset()
	notifierOnce = sync.Once{}

	expectedConfig := &config.Config{
		Notifier: config.NotifierConfig{
			Type:     "telegram",
			TgChatID: "-1001234567890",
		},
	}

	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

	data, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = Get()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGet_ErrorEmptyChatID(t *testing.T) {
	defer config.Reset()
	notifierOnce = sync.Once{}

	expectedConfig := &config.Config{
		Notifier: config.NotifierConfig{
			Type:    "telegram",
			TgToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		},
	}

	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

	data, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = Get()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGet_ErrorUnknownType(t *testing.T) {
	defer config.Reset()
	notifierOnce = sync.Once{}

	expectedConfig := &config.Config{
		Notifier: config.NotifierConfig{
			Type: "unknown",
		},
	}

	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

	data, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = Get()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGet_ErrorLoadConfig(t *testing.T) {
	defer config.Reset()
	notifierOnce = sync.Once{}

	os.Setenv("PINGATUS_CONFIG_PATH", "")

	_, err := Get()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
