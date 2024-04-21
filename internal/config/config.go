package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents config file.
type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Server   Server   `yaml:"server"`
}

// Postgres represents postgres config.
type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Server represents server config.
type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Load loads config from configPath.
func Load(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
