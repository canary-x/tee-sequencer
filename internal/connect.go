package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ConnectHandlerOptions struct {
	Timeout time.Duration `envconfig:"CONNECT_HANDLER_TIMEOUT" default:"3s"`

	ReadTimeout time.Duration `envconfig:"CONNECT_READ_TIMEOUT" default:"3s"`

	// The standard library http.Server.WriteTimeout
	// A zero or negative value means there will be no timeout.
	//
	// https://golang.org/pkg/net/http/#Server.WriteTimeout
	WriteTimeout time.Duration `envconfig:"CONNECT_WRITE_TIMEOUT" default:"3s"`
}

// ConnectServer creates a buf Connect server (https://github.com/bufbuild/connect-go)
// essentially gRPC over HTTP
type ConnectServer struct {
	opt        ConnectHandlerOptions
	mux        *http.ServeMux
	httpServer *http.Server
}

func NewConnectServer(opt ConnectHandlerOptions) *ConnectServer {
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
			defer func() {
				if r := recover(); r != nil {
					resp = nil
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("panic: %v", r))
					log.Printf("Recovering from panic in Connect handler: %+v\n", r)
					log.Printf("Stack trace from panic: %s\n", debug.Stack())
				}
			}()
			resp, err = next(ctx, req)
			if err != nil {
				log.Printf("Error in Connect handler: %s\n", err)
				err = connect.NewError(connect.CodeInternal, err)
			}
			return
		}
	}
}
