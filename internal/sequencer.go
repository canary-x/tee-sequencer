package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/canary-x/tee-sequencer/pkg/api"
	"github.com/mdlayher/vsock"
)

func Run() error {
	var (
		ln  net.Listener
		err error
	)
	cfg, err := ParseConfig()
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if cfg.VSockPort == 0 {
		// fallback to run on http for local testing
		if ln, err = net.Listen("http", fmt.Sprintf(":%d", cfg.HTTPPort)); err != nil {
			return fmt.Errorf("error setting up http listener: %w", err)
		}
	} else {
		// else run on a proper vsock
		if ln, err = vsock.Listen(cfg.VSockPort, nil); err != nil {
			return fmt.Errorf("error setting up vsock listener: %w", err)
		}
	}
	defer ln.Close()

	log.Println("Listening for transactions on vsock...")

	err = http.Serve(ln, http.HandlerFunc(handle))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving http: %w", err)
	}

	log.Println("Server terminated")
	return nil
}

func handle(resp http.ResponseWriter, req *http.Request) {
	var batch api.TransactionBatch
	err := json.NewDecoder(req.Body).Decode(&batch)
	if err != nil {
		http.Error(resp, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	// Don't do any sorting for now, just keep them as they are
	sortedBatch := api.TransactionBatchSorted{
		Transactions: batch.Transactions,
	}

	if err := json.NewEncoder(resp).Encode(sortedBatch); err != nil {
		http.Error(resp, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
