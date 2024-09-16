package internal

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	VSockPort uint32 `envconfig:"VSOCK_PORT"`
	HTTPPort  uint32 `envconfig:"HTTP_PORT"`
}

func (c *Config) Validate() error {
	if c.VSockPort == 0 && c.HTTPPort == 0 {
		return fmt.Errorf("either VSOCK_PORT or HTTP_PORT must be set")
	}
	return nil
}

func ParseConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("parsing from env: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}
