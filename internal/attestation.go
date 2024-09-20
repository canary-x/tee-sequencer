package internal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/canary-x/tee-sequencer/internal/config"
	"github.com/hf/nsm"
	"github.com/hf/nsm/request"
)

// NitroSecurityModule is an interface to The Nitro Secure Module (NSM), which is a hardware security module (HSM) that
// provides various primitives for the secure enclave. This interface exposes all the primitives we need in this app.
type NitroSecurityModule interface {
	// Attest creates a cryptographic attestation of the provided data, to demonstrate that the data was processed
	// within the secure enclave.
	Attest(data []byte) ([]byte, error)
}

// InitSecurityModule enforces initialization of a Nitro Security Module (NSM) session. This only works on AWS EC2
// instances with Nitro installed.
// If the cfg.SecureEnclave option is set to false, it will return a fake implementation, which allows running this
// software on any other machine for development and testing purposes.
func InitSecurityModule(cfg config.Config) (NitroSecurityModule, error) {
	if cfg.SecureEnclave {
		sess, err := nsm.OpenDefaultSession()
		if nil != err {
			return nil, fmt.Errorf("opening NSM session: %w", err)
		}
		// no need to defer closing the session, as it will be closed when the app terminates
		return &secureNSM{session: sess}, nil
	}
	return &fakeNSM{}, nil
}

type secureNSM struct {
	session *nsm.Session
}

func (s *secureNSM) Attest(data []byte) ([]byte, error) {
	res, err := s.session.Send(&request.Attestation{
		UserData: data,
		Nonce:    int64ToBytes(time.Now().UnixMilli()),
	})
	if nil != err {
		return nil, fmt.Errorf("sending attestation request: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("%s", res.Error)
	}
	if nil == res.Attestation || nil == res.Attestation.Document {
		return nil, errors.New("NSM device did not return an attestation")
	}
	return res.Attestation.Document, nil
}

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

type fakeNSM struct{}

func (f *fakeNSM) Attest([]byte) ([]byte, error) {
	return []byte("fake-attestation"), nil
}
