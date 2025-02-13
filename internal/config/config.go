package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var (
	ErrConfigPathNotSet    = errors.New("CONFIG_PATH environment variable not set")
	ErrFailedToOpenFile    = errors.New("config file could not be opened")
	ErrFailedToReadFile    = errors.New("failed to read config file")
	ErrFailedToParseConfig = errors.New("unable to parse config ")
)

type Config struct {
	DBPath string `yaml:"db_path"`
	Secret string `yaml:"secret"`
}

func Init() (*Config, error) {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		return nil, ErrConfigPathNotSet
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, ErrFailedToOpenFile
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrFailedToReadFile
	}

	var cfg *Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, ErrFailedToParseConfig
	}

	return cfg, nil
}
