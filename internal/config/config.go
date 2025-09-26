package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Enabled    *bool  `yaml:"enabled"`
	IPSource   Source `yaml:"ip_source"`
	PortSource Source `yaml:"port_source"`
}

type Source struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

func LoadFromFile(filePath string) (Config, error) {
	var conf Config

	content, err := os.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader, yaml.DisallowUnknownField())

	if err := decoder.Decode(&conf); err != nil {
		return conf, fmt.Errorf("error parsing config yaml: %w", err)
	}

	return conf, nil
}
