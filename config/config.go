package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug     bool              `yaml:"debug,omitempty"`
	MongoURI  string            `yaml:"mongouri"`
	HTTPPoint []HTTPpointConfig `yaml:"httppoint,omitempty"`
}

type HTTPpointConfig struct {
	Name     string        `yaml:"name"`
	Url      string        `yaml:"url"`
	Status   int           `yaml:"status"`
	Timeout  time.Duration `yaml:"timeout"`
	Interval time.Duration `yaml:"interval"`
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

		fmt.Println("Config loaded.")
	})

	return config, configError
}
