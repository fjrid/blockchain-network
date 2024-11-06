package mpt

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

type (
	ExtensionNode struct {
		Path  []byte
		Child Node
	}
)

func (n *ExtensionNode) Hash() []byte {
	serialized, err := serialize(n)
	if err != nil {
		log.Fatalf("failed to hash node: %+v", err)
	}

	return crypto.Keccak256(serialized)
}

func (n *ExtensionNode) Raw() []interface{} {
	result := make([]interface{}, 2)
	result[0] = nibbleToByte(setPrefix(n.Path, false))

	serialized, err := serialize(n.Child)
	if err != nil {
		log.Fatalf("failed to serialize data: %+v", err)
	}

	if len(serialized) >= 32 {
		result[1] = []byte(n.Child.Hash())
	} else {
		result[1] = n.Child.Raw()
	}

	return result
}
