package internal

import (
	"context"
	"errors"
	"fmt"

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

type SequencerServiceHandler struct {
	nsm NitroSecurityModule
}

func NewSequencerServiceHandler(nsm NitroSecurityModule) *SequencerServiceHandler {
	return &SequencerServiceHandler{nsm: nsm}
}

var _ blockchainv1connect.SequencerServiceHandler = (*SequencerServiceHandler)(nil)

func (h *SequencerServiceHandler) Shuffle(
	_ context.Context, req *connect.Request[v1.ShuffleRequest],
) (*connect.Response[v1.ShuffleResponse], error) {
	log := logger.Instance()
	log.Info("Handling shuffle request")

	sequenceOut := req.Msg.Transactions // do not shuffle for now, the algo will be implemented subsequently

	attestation, err := h.attestSequence(sequenceOut)
	if err != nil {
		var clientErr errInvalidTransactionHash
		if errors.As(err, &clientErr) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		// server error
		return nil, err
	}
	log.Debug("Attestation created")

	return connect.NewResponse(&v1.ShuffleResponse{
		Transactions: req.Msg.Transactions,
		Attestation:  attestation,
	}), nil
}

// This method produces a Nitro attestation for the sequence of transactions, by calculating the keccak256 hash of
// all transaction hashes and feeding that as UserData for the attestation.
// Returns errInvalidTransactionHash if any transaction hash is not 32 bytes.
func (h *SequencerServiceHandler) attestSequence(sequenceOut []*v1.Transaction) ([]byte, error) {
	getTxHash := func(tx *v1.Transaction) []byte {
		return tx.TxHash
	}
	// calculate the keccak-256 hash of all sequence tx hashes
	hash := objectKeccak256(sequenceOut, getTxHash)
	// attest the hash
	return h.nsm.Attest(hash)
}

type errInvalidTransactionHash struct {
	TxIdx int
}

func (e errInvalidTransactionHash) Error() string {
	return fmt.Sprintf("transaction at index %d has an invalid hash: should be 32 bytes", e.TxIdx)
}
