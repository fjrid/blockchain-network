package mpt

import (
	"github.com/ethereum/go-ethereum/rlp"
)

type (
	Node interface {
		Hash() []byte
		Raw() []interface{}
	}
)

func serialize(n Node) ([]byte, error) {
	var raw interface{}
	if n == nil {
		raw = [][]byte{}
	} else {
		raw = n.Raw()
	}

	rlp, err := rlp.EncodeToBytes(raw)
	if err != nil {
		return nil, err
	}

	return rlp, nil
}
