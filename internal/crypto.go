package internal

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func newKeccakState() keccakState {
	return sha3.NewLegacyKeccak256().(keccakState)
}

func objectKeccak256[T any](objects []T, getBytes func(obj T) []byte) []byte {
	b := make([]byte, 32)
	d := newKeccakState()
	for _, object := range objects {
		payload := getBytes(object)
		_, _ = d.Write(payload)
	}
	_, _ = d.Read(b)
	return b
}
