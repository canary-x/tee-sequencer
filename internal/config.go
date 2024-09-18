package internal

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	VSockPort uint32 `envconfig:"VSOCK_PORT" default:"8080"`
	Connect   ConnectHandlerOptions
}

func (c *Config) Validate() error {
	return nil // unnecessary for now, but can implement something in the future
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
