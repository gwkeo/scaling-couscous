package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	DBPath string `yaml:"db_path"`
	Secret string `yaml:"secret"`
}

func Init() (*Config, error) {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		return nil, errors.New("CONFIG_PATH environment variable not set")
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, errors.New("failed to open config file: " + err.Error())
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	var cfg *Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.New("unable to parse config: " + err.Error())
	}

	return cfg, nil
}
