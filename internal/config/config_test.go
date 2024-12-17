package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestLoad_ConfigPathNotSet(t *testing.T) {
	os.Unsetenv("PINGATUS_CONFIG_PATH")

	_, err := Load()
	if err != ErrPathNotSet {
		t.Fatalf("expected error %v, got %v", ErrPathNotSet, err)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	os.Setenv("PINGATUS_CONFIG_PATH", "/non/existent/path.yaml")
	defer os.Unsetenv("PINGATUS_CONFIG_PATH")

	_, err := Load()
	if err == nil || !errors.Is(err, ErrReadingFile) {
		t.Fatalf("expected %v, got %v", ErrReadingFile, err)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid_config_*.yaml")
	if err != nil {
		t.Fatalf("can not create temporary invalid file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidYAML := "invalid_yaml: - name: value"
	tmpFile.WriteString(invalidYAML)
	tmpFile.Close()

	os.Setenv("PINGATUS_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("PINGATUS_CONFIG_PATH")

	_, err = Load()
	if err == nil || !errors.Is(err, ErrParsingFile) {
		t.Fatalf("expected %v, got %v", ErrParsingFile, err)
	}
}

func TestLoad_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "valid_config_*.yaml")
	if err != nil {
		t.Fatalf("can not create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validConfig := Config{
		Endpoints: []EndpointConfig{
			{
				Name:        "test-endpoint",
				Type:        "http",
				Address:     "http://example.com",
				Status:      200,
				PacketCount: 5,
				Timeout:     10 * time.Second,
				Interval:    5 * time.Second,
			},
		},
	}

	data, err := yaml.Marshal(validConfig)
	if err != nil {
		t.Fatalf("can not serialize YAML: %v", err)
	}

	tmpFile.Write(data)
	tmpFile.Close()

	os.Setenv("PINGATUS_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("PINGATUS_CONFIG_PATH")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(cfg.Endpoints))
	}

	endpoint := cfg.Endpoints[0]
	if endpoint.Name != "test-endpoint" || endpoint.Type != "http" {
		t.Errorf("eunexpected endpoint data: %+v", endpoint)
	}
}
