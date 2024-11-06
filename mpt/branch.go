package mpt

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

type (
	BranchNode struct {
		Children [16]Node
		Value    []byte
	}
)

func (n *BranchNode) Hash() []byte {
	serialized, err := serialize(n)
	if err != nil {
		log.Fatalf("failed to hash node: %+v", err)
	}

	return crypto.Keccak256(serialized)
}

func (n *BranchNode) Raw() []interface{} {
	result := make([]interface{}, 17)

	for i := 0; i < len(n.Children); i++ {
		if n.Children[i] == nil {
			result[i] = []byte{}
			continue
		}

		serialized, err := serialize(n.Children[i])
		if err != nil {
			log.Fatalf("failed to serialize data: %+v", err)
		}

		if len(serialized) >= 32 {
			result[i] = n.Children[i].Hash()
		} else {
			result[i] = n.Children[i].Raw()
		}
	}

	result[16] = n.Value

	return result
}
