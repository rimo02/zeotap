package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"time"
)

type Threshold struct {
	MaximumTemp float64 `yaml:"maximumTemp"`
	Breach      int     `yaml:"breach"`
}

type City struct {
	Name         string    `yaml:"name"`
	TempThreshold Threshold `yaml:"tempThreshold"`
}

type Config struct {
	Interval  time.Duration `yaml:"interval"`
	Cities    []City   `yaml:"cities"`
	TempUnit  string   `yaml:"tempUnit"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	return cfg
}

