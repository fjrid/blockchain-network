package transaction

import (
	"crypto/sha256"
	"errors"
)

type (
	MerkleNode struct {
		Left  *MerkleNode
		Right *MerkleNode
		Hash  []byte
	}

	MerkleTree struct {
		Root *MerkleNode
	}
)

func newMerkleNode(left, right *MerkleNode, hash []byte) *MerkleNode {
	if left == nil && right == nil {
		return &MerkleNode{
			Hash: hash,
		}
	}

	node := &MerkleNode{
		Left:  left,
		Right: right,
	}

	secureHash := sha256.New()
	secureHash.Write(left.Hash)
	secureHash.Write(right.Hash)
	node.Hash = secureHash.Sum(nil)

	return node
}

func MerkleHashTransactions(transactions []*Transaction) ([]byte, error) {
	if len(transactions) == 0 {
		return nil, errors.New("invalid transactions")
	}

	nodes := make([]*MerkleNode, 0)

	for _, transaction := range transactions {
		nodes = append(nodes, newMerkleNode(nil, nil, transaction.Hash()))
	}

	for len(nodes) > 1 {
		tempNodes := make([]*MerkleNode, 0)

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				tempNodes = append(tempNodes, newMerkleNode(nodes[i], nodes[i+1], nil))
			} else {
				tempNodes = append(tempNodes, nodes[i])
			}
		}

		nodes = tempNodes
	}

	return nodes[0].Hash, nil
}
