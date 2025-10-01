package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type SourceType string

const (
	SourceTypeHTTP   SourceType = "http"
	SourceTypeFile   SourceType = "file"
	SourceTypeStatic SourceType = "static"
)

type Config struct {
	Services map[string]Service `yaml:"services" validate:"required,dive"`
}

type Service struct {
	Enabled    *bool         `yaml:"enabled" default:"true" validate:"required"`
	Interval   time.Duration `yaml:"interval" validate:"required,min=30s"`
	IPSource   Source        `yaml:"ip_source" validate:"required"`
	PortSource Source        `yaml:"port_source" validate:"required"`
}

type Source struct {
	Type       SourceType    `yaml:"type" validate:"required,oneof=http file static"`
	Value      string        `yaml:"value" validate:"required_if=Type static"`
	URL        string        `yaml:"url" validate:"required_if=Type http"`
	JSONPath   string        `yaml:"json_path"`
	Timeout    time.Duration `yaml:"timeout" default:"10s"`
	Path       string        `yaml:"path" validate:"required_if=Type file"`
	Pattern    string        `yaml:"pattern"`
	MatchGroup *int          `yaml:"match_group" default:"1" validate:"omitempty,min=0"`
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

	if err := defaults.Set(&conf); err != nil {
		return conf, fmt.Errorf("error setting config defaults: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(conf); err != nil {
		// TODO: Pretty errors
		return conf, err
	}

	return conf, nil
}
