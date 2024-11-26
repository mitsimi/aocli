package config

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Session   string `json:"session" yaml:"session" toml:"session"`
	Year      int    `json:"year" yaml:"year" toml:"year"`
	Structure string `json:"structure" yaml:"structure" toml:"structure"`
}

func ParseJSON(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseYAML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseTOML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if _, err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
