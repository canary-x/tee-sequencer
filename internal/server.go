package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/canary-x/tee-sequencer/pkg/api"
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

	log.Println("Listening for transactions on vsock...")

	err = http.Serve(ln, http.HandlerFunc(handle))
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
}
