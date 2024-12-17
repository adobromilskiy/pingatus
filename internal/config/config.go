package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrPathNotSet  = errors.New("config: PINGATUS_CONFIG_PATH is not set")
	ErrReadingFile = errors.New("config: error reading yaml file")
	ErrParsingFile = errors.New("config: error parsing yaml file")
)

type Config struct {
	Logger    LoggerConfig     `yaml:"logger"`
	Endpoints []EndpointConfig `yaml:"endpoints,omitempty"`
}

type EndpointConfig struct {
	Name        string        `yaml:"name"`
	Type        string        `yaml:"type"`
	Address     string        `yaml:"address"`
	Status      int           `yaml:"status"`
	PacketCount int           `yaml:"packetcount"`
	Timeout     time.Duration `yaml:"timeout"`
	Interval    time.Duration `yaml:"interval"`
}

type LoggerConfig struct {
	IsJSON bool   `yaml:"json,omitempty"`
	Level  string `yaml:"level,omitempty"`
}

func Load() (*Config, error) {
	path := os.Getenv("PINGATUS_CONFIG_PATH")
	if len(path) == 0 {
		return nil, ErrPathNotSet
	}

	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadingFile, err)
	}

	var cfg Config

	err = yaml.Unmarshal(bts, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParsingFile, err)
	}

	return &cfg, nil
}
