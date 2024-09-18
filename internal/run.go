package internal

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1/blockchainv1connect"
	"github.com/mdlayher/vsock"
)

func Run() error {
	cfg, err := ParseConfig()
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	ln, err := listen(cfg)
	if err != nil {
		return fmt.Errorf("listening on socket: %w", err)
	}
	defer ln.Close()

	log.Println("Listening for transactions...")

	interceptors := connect.WithInterceptors(ConnectErrorInterceptor())
	srv := NewConnectServer(cfg.Connect).
		WithHandler(blockchainv1connect.NewPingServiceHandler(NewPingServiceHandler(), interceptors)).
		WithHandler(blockchainv1connect.NewSequencerServiceHandler(NewSequencerServiceHandler(), interceptors))

	err = srv.Serve(ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving http: %w", err)
	}

	log.Println("Server terminated")
	return nil
}

func listen(cfg Config) (net.Listener, error) {
	ln, err := vsock.Listen(cfg.VSockPort, nil)
	if err != nil && strings.Contains(err.Error(), "vsock: not implemented") {
		log.Println("OS does not support vsock: falling back to regular TCP socket")
		return net.Listen("tcp", fmt.Sprintf(":%d", cfg.VSockPort))
	}
	return ln, err
}
