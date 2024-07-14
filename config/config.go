package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug     bool             `yaml:"debug,omitempty"`
	MongoURI  string           `yaml:"mongouri"`
	EndPoints []EndpointConfig `yaml:"endpoints,omitempty"`
	WEBAPI    WEBAPIConfig     `yaml:"webapi"`
	Notifier  NotifierConfig   `yaml:"notifier,omitempty"`
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

type WEBAPIConfig struct {
	ListenAddr string `yaml:"listenaddr"`
	AssetsDir  string `yaml:"assetsdir"`
}

type NotifierConfig struct {
	Type     string `yaml:"type"`
	TgToken  string `yaml:"tgtoken,omitempty"`
	TgChatID string `yaml:"tgchatid,omitempty"`
}

var (
	configOnce  sync.Once
	config      *Config
	configError error
)

func Load() (*Config, error) {
	configOnce.Do(func() {
		configPath := os.Getenv("PINGATUS_CONFIG_PATH")
		if len(configPath) == 0 {
			configError = fmt.Errorf("PINGATUS_CONFIG_PATH is not set")
			return
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			configError = fmt.Errorf("Error reading YAML file: %v", err)
			return
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			configError = fmt.Errorf("Error parsing YAML file: %v", err)
			return
		}

		log.Println("[INFO] config loaded")
	})

	return config, configError
}
