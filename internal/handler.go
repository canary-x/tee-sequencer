package internal

import (
	"context"

	"github.com/bufbuild/connect-go"
	v1 "github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1"
	"github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1/blockchainv1connect"
	"github.com/canary-x/tee-sequencer/internal/logger"
)

type PingServiceHandler struct{}

func NewPingServiceHandler() *PingServiceHandler {
	return &PingServiceHandler{}
}

var _ blockchainv1connect.PingServiceHandler = (*PingServiceHandler)(nil)

func (h *PingServiceHandler) Ping(
	context.Context, *connect.Request[v1.PingRequest],
) (*connect.Response[v1.PingResponse], error) {
	return connect.NewResponse(&v1.PingResponse{
		Message: "pong",
	}), nil
}

type SequencerServiceHandler struct{}

func NewSequencerServiceHandler() *SequencerServiceHandler {
	return &SequencerServiceHandler{}
}

var _ blockchainv1connect.SequencerServiceHandler = (*SequencerServiceHandler)(nil)

func (h *SequencerServiceHandler) Shuffle(
	_ context.Context, req *connect.Request[v1.ShuffleRequest],
) (*connect.Response[v1.ShuffleResponse], error) {
	logger.Instance().Debug("Handling shuffle request")
	return connect.NewResponse(&v1.ShuffleResponse{
		Transactions: req.Msg.Transactions,
		Signature:    nil,
	}), nil
}
