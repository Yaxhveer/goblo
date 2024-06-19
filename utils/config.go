package utils

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type LBStrategy int

const (
	RoundRobin LBStrategy = iota
	LeastConnected
	Random
)

func GetLBStrategy(strategy string) LBStrategy {
	switch strategy {
	case "least-connection":
		return LeastConnected
	case "random":
		return Random
	default:
		return RoundRobin
	}
}

type Config struct {
	Port            int      `json:"port" yaml:"port"`
	Backends        []string `json:"backends" yaml:"backends"`
	Strategy        string   `json:"strategy" yaml:"strategy"`
}

const MAX_LB_ATTEMPTS int = 3

func GetLBConfig(logger *zap.Logger) (*Config, error) {
	var config Config
	configFile, err := os.ReadFile("config.yaml")
	if err == os.ErrNotExist {
		logger.Warn("config file not found")
		return &Config{
			Port: 3333,
			Strategy: "least-connection",
			Backends: []string{
				"http://localhost:3334",
				"http://localhost:3335",
				"http://localhost:3336",
				"http://localhost:3337",
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	if len(config.Backends) == 0 {
		return nil, errors.New("backend hosts expected, none provided")
	}

	if config.Port == 0 {
		return nil, errors.New("load balancer port not found")
	}

	return &config, nil
}
