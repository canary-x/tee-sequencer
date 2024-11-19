package internal_test

import (
	"context"
	"testing"

	"github.com/bufbuild/connect-go"
	v1 "github.com/canary-x/tee-sequencer/gen/proto/go/blockchain/v1"
	"github.com/canary-x/tee-sequencer/internal"
	"github.com/canary-x/tee-sequencer/internal/logger"
	"github.com/stretchr/testify/suite"
)

type SequencerServiceHandlerTest struct {
	suite.Suite
	handler *internal.SequencerServiceHandler
}

func TestSequencerServiceHandlerTest(t *testing.T) {
	suite.Run(t, new(SequencerServiceHandlerTest))
}

func (s *SequencerServiceHandlerTest) SetupSuite() {
	logger.InitForTests()
	s.handler = internal.NewSequencerServiceHandler(&internal.FakeNSM{})
}

func (s *SequencerServiceHandlerTest) TestShuffle_Ok() {
	tx1 := &v1.Transaction{
		TxHash:  []byte("11111111111111111111111111111111"),
		Account: []byte("1"),
		Nonce:   []byte("1"),
	}
	tx2 := &v1.Transaction{
		TxHash:  []byte("22222222222222222222222222222222"),
		Account: []byte("1"),
		Nonce:   []byte("1"),
	}
	req := connect.NewRequest(&v1.ShuffleRequest{
		Transactions: []*v1.Transaction{
			tx1, tx2,
		},
	})
	resp, err := s.handler.Shuffle(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal("fake-attestation", string(resp.Msg.Attestation))
}

func (s *SequencerServiceHandlerTest) TestShuffle_Validation() {
	tx1 := &v1.Transaction{
		TxHash:  []byte("11111111111111111111111111111111"),
		Account: []byte("1"),
		Nonce:   []byte("1"),
	}
	tx2 := &v1.Transaction{
		TxHash:  []byte("invalid-hash"),
		Account: []byte("1"),
		Nonce:   []byte("1"),
	}
	req := connect.NewRequest(&v1.ShuffleRequest{
		Transactions: []*v1.Transaction{
			tx1, tx2,
		},
	})
	_, err := s.handler.Shuffle(context.Background(), req)
	s.Require().EqualE(err)
}
