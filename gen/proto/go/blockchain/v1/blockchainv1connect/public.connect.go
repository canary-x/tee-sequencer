// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: blockchain/v1/public.proto

package blockchainv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// PingServiceName is the fully-qualified name of the PingService service.
	PingServiceName = "blockchain.v1.PingService"
	// SequencerServiceName is the fully-qualified name of the SequencerService service.
	SequencerServiceName = "blockchain.v1.SequencerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PingServicePingProcedure is the fully-qualified name of the PingService's Ping RPC.
	PingServicePingProcedure = "/blockchain.v1.PingService/Ping"
	// SequencerServiceShuffleProcedure is the fully-qualified name of the SequencerService's Shuffle
	// RPC.
	SequencerServiceShuffleProcedure = "/blockchain.v1.SequencerService/Shuffle"
)

// PingServiceClient is a client for the blockchain.v1.PingService service.
type PingServiceClient interface {
	Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error)
}

// NewPingServiceClient constructs a client for the blockchain.v1.PingService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPingServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) PingServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &pingServiceClient{
		ping: connect_go.NewClient[v1.PingRequest, v1.PingResponse](
			httpClient,
			baseURL+PingServicePingProcedure,
			opts...,
		),
	}
}

// pingServiceClient implements PingServiceClient.
type pingServiceClient struct {
	ping *connect_go.Client[v1.PingRequest, v1.PingResponse]
}

// Ping calls blockchain.v1.PingService.Ping.
func (c *pingServiceClient) Ping(ctx context.Context, req *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error) {
	return c.ping.CallUnary(ctx, req)
}

// PingServiceHandler is an implementation of the blockchain.v1.PingService service.
type PingServiceHandler interface {
	Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error)
}

// NewPingServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPingServiceHandler(svc PingServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	pingServicePingHandler := connect_go.NewUnaryHandler(
		PingServicePingProcedure,
		svc.Ping,
		opts...,
	)
	return "/blockchain.v1.PingService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PingServicePingProcedure:
			pingServicePingHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPingServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPingServiceHandler struct{}

func (UnimplementedPingServiceHandler) Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("blockchain.v1.PingService.Ping is not implemented"))
}

// SequencerServiceClient is a client for the blockchain.v1.SequencerService service.
type SequencerServiceClient interface {
	Shuffle(context.Context, *connect_go.Request[v1.ShuffleRequest]) (*connect_go.Response[v1.ShuffleResponse], error)
}

// NewSequencerServiceClient constructs a client for the blockchain.v1.SequencerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewSequencerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) SequencerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &sequencerServiceClient{
		shuffle: connect_go.NewClient[v1.ShuffleRequest, v1.ShuffleResponse](
			httpClient,
			baseURL+SequencerServiceShuffleProcedure,
			opts...,
		),
	}
}

// sequencerServiceClient implements SequencerServiceClient.
type sequencerServiceClient struct {
	shuffle *connect_go.Client[v1.ShuffleRequest, v1.ShuffleResponse]
}

// Shuffle calls blockchain.v1.SequencerService.Shuffle.
func (c *sequencerServiceClient) Shuffle(ctx context.Context, req *connect_go.Request[v1.ShuffleRequest]) (*connect_go.Response[v1.ShuffleResponse], error) {
	return c.shuffle.CallUnary(ctx, req)
}

// SequencerServiceHandler is an implementation of the blockchain.v1.SequencerService service.
type SequencerServiceHandler interface {
	Shuffle(context.Context, *connect_go.Request[v1.ShuffleRequest]) (*connect_go.Response[v1.ShuffleResponse], error)
}

// NewSequencerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewSequencerServiceHandler(svc SequencerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	sequencerServiceShuffleHandler := connect_go.NewUnaryHandler(
		SequencerServiceShuffleProcedure,
		svc.Shuffle,
		opts...,
	)
	return "/blockchain.v1.SequencerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case SequencerServiceShuffleProcedure:
			sequencerServiceShuffleHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedSequencerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedSequencerServiceHandler struct{}

func (UnimplementedSequencerServiceHandler) Shuffle(context.Context, *connect_go.Request[v1.ShuffleRequest]) (*connect_go.Response[v1.ShuffleResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("blockchain.v1.SequencerService.Shuffle is not implemented"))
}
