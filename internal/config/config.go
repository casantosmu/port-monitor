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

type CheckType string

const (
	CheckTypeHTTP CheckType = "http"
)

type Config struct {
	Services map[string]Service `yaml:"services" validate:"required,dive"`
}

type Service struct {
	Enabled           *bool         `yaml:"enabled" default:"true" validate:"required"`
	Interval          time.Duration `yaml:"interval" validate:"required,min=30s"`
	IPSource          Source        `yaml:"ip_source" validate:"required"`
	PortSource        Source        `yaml:"port_source" validate:"required"`
	ConnectivityCheck Check         `yaml:"connectivity_check" validate:"required"`
}

type Source struct {
	Type       SourceType        `yaml:"type" validate:"required,oneof=http file static"`
	Value      string            `yaml:"value" validate:"required_if=Type static"`
	URL        string            `yaml:"url" validate:"required_if=Type http"`
	Headers    map[string]string `yaml:"headers"`
	Method     string            `yaml:"method" default:"GET" validate:"oneof=GET POST"`
	Body       string            `yaml:"body"`
	Proxy      string            `yaml:"proxy"`
	JSONPath   string            `yaml:"json_path"`
	Timeout    time.Duration     `yaml:"timeout" default:"10s"`
	Path       string            `yaml:"path" validate:"required_if=Type file"`
	Pattern    string            `yaml:"pattern"`
	MatchGroup *int              `yaml:"match_group" default:"1" validate:"omitempty,min=0"`
}

type Check struct {
	Type           CheckType         `yaml:"type" validate:"required,oneof=http"`
	URL            string            `yaml:"url" validate:"required_if=Type http"`
	Headers        map[string]string `yaml:"headers"`
	Method         string            `yaml:"method" default:"GET" validate:"oneof=GET POST"`
	Body           string            `yaml:"body"`
	Proxy          string            `yaml:"proxy"`
	SuccessPattern string            `yaml:"success_pattern" validate:"required_if=Type http"`
	Timeout        time.Duration     `yaml:"timeout" default:"10s"`
}

func LoadFromFile(filePath string) (Config, error) {
	var conf Config

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	content := expandVars(raw)

	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader, yaml.DisallowUnknownField())

	if err := decoder.Decode(&conf); err != nil {
		rawReader := bytes.NewReader(raw)
		rawDecoder := yaml.NewDecoder(rawReader, yaml.DisallowUnknownField())

		if err := rawDecoder.Decode(&conf); err != nil {
			return conf, fmt.Errorf("error parsing config yaml: %w", err)
		}

		// Security: If raw config is valid, the env vars broke the YAML.
		// Return generic error to avoid leaking secrets in the logs.
		return conf, fmt.Errorf("error parsing config yaml after env var expansion (error details hidden for security)")
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
