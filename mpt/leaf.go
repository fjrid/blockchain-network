package mpt

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

type (
	LeafNode struct {
		Key   []byte
		Value []byte
	}
)

func (n *LeafNode) Hash() []byte {
	serialized, err := serialize(n)
	if err != nil {
		log.Fatalf("failed to hash node: %+v", err)
	}

	return crypto.Keccak256(serialized)
}

func (n *LeafNode) Raw() []interface{} {
	return []interface{}{
		nibbleToByte(setPrefix(n.Key, true)),
		n.Value,
	}
}
