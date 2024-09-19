package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/bufbuild/connect-go"
	"github.com/canary-x/tee-sequencer/internal/config"
	"github.com/canary-x/tee-sequencer/internal/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// ConnectServer creates a buf Connect server (https://github.com/bufbuild/connect-go)
// essentially gRPC over HTTP
type ConnectServer struct {
	opt        config.ConnectHandlerOptions
	mux        *http.ServeMux
	httpServer *http.Server
}

func NewConnectServer(opt config.ConnectHandlerOptions) *ConnectServer {
	return &ConnectServer{
		opt: opt,
		mux: http.NewServeMux(),
	}
}

func (s *ConnectServer) WithHandler(pattern string, handler http.Handler) *ConnectServer {
	s.mux.Handle(pattern, handler)
	return s
}

func (s *ConnectServer) Serve(ln net.Listener) error {
	s.httpServer = &http.Server{
		ReadHeaderTimeout: s.opt.ReadTimeout,
		WriteTimeout:      s.opt.WriteTimeout,
		Handler:           h2c.NewHandler(s.mux, &http2.Server{}), // Use h2c, so we can serve HTTP/2 without TLS.
	}
	err := s.httpServer.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// ConnectErrorInterceptor recovers from panics inside of handlers and prints them as an error log with a stack trace.
// This is quite important as we definitely don't want a panic to stop our enclave.
// It also handles regular errors and ensures the error message is both logged and included in the response.
// Without this interceptor, a client would just see "internal" as the error message, which is fine in traditional
// client-server set-ups, but not in our case, as this enclave is supposed to be embedded, and the parent service needs
// to be able to expose errors to facilitate debugging.
func ConnectErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			log := logger.Instance()
			defer func() {
				if r := recover(); r != nil {
					resp = nil
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("panic: %v", r))
					log.Error("Recovering from panic in Connect handler: %+v\n", r)
					log.Error("Stack trace from panic: %s\n", debug.Stack())
				}
			}()
			resp, err = next(ctx, req)
			if err != nil {
				log.Error("Error in Connect handler: %s\n", err)
				err = connect.NewError(connect.CodeInternal, err)
			}
			return
		}
	}
}
