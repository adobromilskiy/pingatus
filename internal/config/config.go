package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	errPathNotSet  = errors.New("config: PINGATUS_CONFIG_PATH is not set")
	errReadingFile = errors.New("config: error reading yaml file")
	errParsingFile = errors.New("config: error parsing yaml file")
)

type Config struct {
	DBDSN     string           `yaml:"dbdsn"`
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
		return nil, errPathNotSet
	}

	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errReadingFile, err)
	}

	var cfg Config

	err = yaml.Unmarshal(bts, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errParsingFile, err)
	}

	return &cfg, nil
}
