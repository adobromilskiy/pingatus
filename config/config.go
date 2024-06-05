package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug    bool   `yaml:"debug,omitempty"`
	MongoURI string `yaml:"mongouri"`
}

var (
	configOnce  sync.Once
	config      *Config
	configError error
)

func Load() (*Config, error) {
	configOnce.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if len(configPath) == 0 {
			configError = fmt.Errorf("CONFIG_PATH is not set")
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
	})

	return config, configError
}
