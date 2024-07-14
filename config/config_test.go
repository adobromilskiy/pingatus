package config

import (
	"os"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoad(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

	expectedConfig := &Config{
		Debug:    true,
		MongoURI: "mongodb://localhost:27017",
		EndPoints: []EndpointConfig{
			{
				Name:     "test",
				Address:  "http://localhost",
				Status:   200,
				Timeout:  1,
				Interval: 1,
			},
		},
		WEBAPI: WEBAPIConfig{
			ListenAddr: "localhost:8080",
			AssetsDir:  "/assets",
		},
	}

	data, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Debug != expectedConfig.Debug {
		t.Errorf("Expected %v, got %v", expectedConfig, config)
	}

	configOnce = sync.Once{}

	if _, err := tempFile.Write([]byte("invalid yaml")); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestLoad_NoConfig(t *testing.T) {
	configOnce = sync.Once{}

	os.Setenv("PINGATUS_CONFIG_PATH", "non-existent-file")

	_, err := Load()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestLoad_EmptyEnvConfigPath(t *testing.T) {
	configOnce = sync.Once{}

	os.Setenv("PINGATUS_CONFIG_PATH", "")

	_, err := Load()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
