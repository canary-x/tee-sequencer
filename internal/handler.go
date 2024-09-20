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

	sequenceIn := req.Msg.Transactions
	sequenceOut := sequenceIn // do not shuffle for now, the algo will be implemented subsequently

	attestation, err := AttestSequence(sequenceIn, sequenceOut, h.nsm)
	if err != nil {
		if errors.Is(err, errInvalidTransactionHash{}) { // TODO does this work or do we need to use errors.As?
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}
	log.Debug("Attestation created")

	return connect.NewResponse(&v1.ShuffleResponse{
		Transactions: req.Msg.Transactions,
		Attestation:  attestation,
	}), nil
}

func AttestSequence(sequenceIn, sequenceOut []*v1.Transaction, nsm NitroSecurityModule) ([]byte, error) {
	combined := append(sequenceIn, sequenceOut...)
	hashes, errIdx := flattenTransactions(combined)
	if errIdx != -1 {
		return nil, errInvalidTransactionHash{TxIdx: errIdx}
	}
	return nsm.Attest(hashes)
}

// flattenTransactions returns a byte slice containing the hashes of all transactions in the input slice.
// It returns the slice and -1 if the operation completed successfully.
// If any transaction has an invalid hash (not 32 bytes), it returns nil and the index of the invalid transaction.
func flattenTransactions(txs []*v1.Transaction) ([]byte, int) {
	result := make([]byte, 0, len(txs)*32) // each transaction hash is 32 bytes
	var pos int
	for i, tx := range txs {
		if len(tx.TxHash) != 32 {
			return nil, i
		}
		pos += copy(result[pos:], tx.TxHash)
	}
	return result, -1
}

type errInvalidTransactionHash struct {
	TxIdx int
}

func (e errInvalidTransactionHash) Error() string {
	return fmt.Sprintf("transaction at index %d has an invalid hash: should be 32 bytes", e.TxIdx)
}
