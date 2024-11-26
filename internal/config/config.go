package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Session   string `json:"session" yaml:"session" toml:"session"`
	Year      int    `json:"year" yaml:"year" toml:"year"`
	Structure string `json:"structure" yaml:"structure" toml:"structure"`
}

func Parse(path string) (*Config, error) {
	switch ext := filepath.Ext(path); ext {
	case ".json":
		return ParseJSON(path)
	case ".yaml", ".yml":
		return ParseYAML(path)
	case ".toml":
		return ParseTOML(path)
	default:
		return nil, fmt.Errorf("unsupported config file extension: %s", ext)
	}
}

func ParseJSON(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
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

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
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

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Write(path string) error {
	switch ext := filepath.Ext(path); ext {
	case ".json":
		return c.WriteJSON(path)
	case ".yaml", ".yml":
		return c.WriteYAML(path)
	case ".toml":
		return c.WriteTOML(path)
	default:
		return fmt.Errorf("unsupported config file extension: %s", ext)
	}
}

func (c *Config) WriteJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) WriteYAML(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) WriteTOML(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
