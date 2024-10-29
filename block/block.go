package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(fmt.Sprintf("%d", b.Timestamp))
	headers := bytes.Join([][]byte{b.Data, b.PrevBlockHash, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: preBlockHash,
	}
	block.SetHash()

	return block
}

func GenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"), []byte{})
}