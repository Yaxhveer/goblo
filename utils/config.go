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
	Port     int      `json:"port" yaml:"port"`
	Backends []string `json:"backends" yaml:"backends"`
	Strategy string   `json:"strategy" yaml:"strategy"`
}

func GetLBConfig(configFileName string, logger *zap.Logger) (*Config, error) {
	var config Config
	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		logger.Warn("config file error, starting ", zap.Error(err))
		return &Config{
			Port:     3330,
			Strategy: "least-connection",
			Backends: []string{
				"http://localhost:3331",
				"http://localhost:3332",
				"http://localhost:3333",
				"http://localhost:3334",
			},
		}, nil
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
