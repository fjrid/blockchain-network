package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

type (
	Node interface {
		Hash() string
	}

	LeafNode struct {
		Key   []byte
		Value []byte
	}

	BranchNode struct {
		Children [16]Node
		Value    []byte
	}

	ExtensionNode struct {
		Path  []byte
		Child Node
	}

	MerklePatriciaTrie struct {
		Root Node
	}
)

func NewMerklePatriciaTrie() *MerklePatriciaTrie {
	return &MerklePatriciaTrie{}
}

func (n *LeafNode) Hash() string {
	hasher := sha256.New()
	hasher.Write(n.Key)
	hasher.Write(n.Value)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (n *BranchNode) Hash() string {
	hasher := sha256.New()
	for _, child := range n.Children {
		if child != nil {
			hasher.Write([]byte(child.Hash()))
		}
	}
	hasher.Write(n.Value)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (n *ExtensionNode) Hash() string {
	hasher := sha256.New()
	hasher.Write(n.Path)
	hasher.Write([]byte(n.Child.Hash()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (t *MerklePatriciaTrie) Insert(key, value []byte) {
	t.Root = insertNode(t.Root, getNibbleKey(key), value)
}

func insertNode(node Node, key, value []byte) Node {
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
		n.Children[index] = insertNode(n.Children[index], key[1:], value)
		return n
	case *ExtensionNode:
		commonPrefix, remainingKey1, remainingKey2 := getLongestCommonPrefix(n.Path, key)
		if bytes.Equal(n.Path, commonPrefix) {
			n.Child = insertNode(n.Child, remainingKey2, value)
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

func getNibbleKey(keys []byte) []byte {
	result := make([]byte, 0)

	for _, key := range keys {
		result = append(result, key>>4)
		result = append(result, key&0xf)
	}

	return result
}

func getLongestCommonPrefix(key1, key2 []byte) (commonPrefix, remainingKey1, remainingKey2 []byte) {
	i := 0
	for i < len(key1) && i < len(key2) && key1[i] == key2[i] {
		i++
	}
	return key1[:i], key1[i:], key2[i:]
}
