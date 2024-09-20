package internal

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	v1 "github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1/blockchainv1connect"
	"github.com/canary-x/tee-sequencer/internal/config"
	"github.com/canary-x/tee-sequencer/internal/logger"
	"github.com/mdlayher/vsock"
)

func Run() error {
	cfg, err := config.Parse()
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	log := logger.Init(cfg)

	ln, err := listen(cfg, log)
	if err != nil {
		return fmt.Errorf("listening on socket: %w", err)
	}
	defer ln.Close()

	securityModule, err := InitSecurityModule(cfg)
	if err != nil {
		return fmt.Errorf("initializing security module: %w", err)
	}

	log.Info("Listening for transactions...")

	interceptors := connect.WithInterceptors(ConnectErrorInterceptor())
	srv := NewConnectServer(cfg.Connect).
		WithHandler(v1.NewPingServiceHandler(NewPingServiceHandler(), interceptors)).
		WithHandler(v1.NewSequencerServiceHandler(NewSequencerServiceHandler(securityModule), interceptors))

	err = srv.Serve(ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving http: %w", err)
	}

	log.Info("Server terminated")
	return nil
}

func listen(cfg config.Config, log logger.Logger) (net.Listener, error) {
	ln, err := vsock.Listen(cfg.VSockPort, nil)
	if err != nil && strings.Contains(err.Error(), "vsock: not implemented") {
		log.Warn("OS does not support vsock: falling back to regular TCP socket")
		return net.Listen("tcp", fmt.Sprintf(":%d", cfg.VSockPort))
	}
	return ln, err
}
