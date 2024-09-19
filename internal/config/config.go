package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	VSockPort    uint32 `envconfig:"VSOCK_PORT" default:"8080"`     // port on which the server listens
	LogVSockPort uint32 `envconfig:"LOG_VSOCK_PORT" default:"8090"` // port to which logs should be streamed
	LogVSockCID  uint32 `envconfig:"LOG_VSOCK_CID" default:"3"`     // CID of the vsock on the host
	Connect      ConnectHandlerOptions
}

func (c *Config) Validate() error {
	return nil // unnecessary for now, but can implement something in the future
}

func Parse() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("parsing from env: %w", err)
	}
	return cfg, nil
}

type ConnectHandlerOptions struct {
	Timeout time.Duration `envconfig:"CONNECT_HANDLER_TIMEOUT" default:"3s"`

	ReadTimeout time.Duration `envconfig:"CONNECT_READ_TIMEOUT" default:"3s"`

	// The standard library http.Server.WriteTimeout
	// A zero or negative value means there will be no timeout.
	//
	// https://golang.org/pkg/net/http/#Server.WriteTimeout
	WriteTimeout time.Duration `envconfig:"CONNECT_WRITE_TIMEOUT" default:"3s"`
}
