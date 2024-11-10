package mpt

import (
	"bytes"
	"log"

	"github.com/fjrid/blockchain-network/db"
)

type (
	MerklePatriciaTrie struct {
		Root Node
		db   *db.DB
	}
)

func NewMerklePatriciaTrie(db *db.DB) *MerklePatriciaTrie {
	return &MerklePatriciaTrie{
		db: db,
	}
}

func (t *MerklePatriciaTrie) Insert(key, value []byte) {
	t.Root = t.insertNode(t.Root, getNibbleKey(key), value)
}

func (t *MerklePatriciaTrie) insertNode(node Node, key, value []byte) Node {
	if node == nil {
		return &LeafNode{Key: key, Value: value}
	}

	switch n := node.(type) {
	case *LeafNode:
		if bytes.Equal(n.Key, key) {
			n.Value = value
			return n
		}

		commonPrefix, remainingKey1, remainingKey2 := getLongestCommonPrefix(n.Key, key)
		branch := &BranchNode{}
		branch.Children[remainingKey1[0]] = &LeafNode{Key: remainingKey1[1:], Value: n.Value}
		branch.Children[remainingKey2[0]] = &LeafNode{Key: remainingKey2[1:], Value: value}

		if len(commonPrefix) > 0 {
			return &ExtensionNode{Path: commonPrefix, Child: branch}
		}
		return branch
	case *BranchNode:
		index := key[0]
		n.Children[index] = t.insertNode(n.Children[index], key[1:], value)
		return n
	case *ExtensionNode:
		commonPrefix, remainingKey1, remainingKey2 := getLongestCommonPrefix(n.Path, key)
		if bytes.Equal(n.Path, commonPrefix) {
			n.Child = t.insertNode(n.Child, remainingKey2, value)
			return n
		}

		branch := BranchNode{}
		branch.Children[remainingKey1[0]] = n
		branch.Children[remainingKey2[0]] = &LeafNode{Key: remainingKey2[1:], Value: value}

		return &ExtensionNode{Path: commonPrefix, Child: &branch}
	default:
		log.Fatal("invalid node type in merkle patricia trie")
	}

	return nil
}

func (t *MerklePatriciaTrie) Store() {
	if t.Root != nil {
		t.storeNode(t.Root)
	}
}

func (t *MerklePatriciaTrie) storeNode(node Node) {
	switch n := node.(type) {
	case *LeafNode:
		serializedData, err := serialize(node)
		if err == nil {
			t.db.Put(n.Hash(), serializedData)
		}
	case *ExtensionNode:
		serializedData, err := serialize(node)
		if err == nil {
			t.db.Put(n.Hash(), serializedData)
		}

		t.storeNode(n.Child)
	case *BranchNode:
		serializedData, err := serialize(node)
		if err == nil {
			t.db.Put(n.Hash(), serializedData)
		}

		for _, child := range n.Children {
			t.storeNode(child)
		}
	}
}
