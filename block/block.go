package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/fjrid/blockchain-network/transaction"
)

type Block struct {
	Timestamp      int64                      `json:"timestamp"`
	Transactions   []*transaction.Transaction `json:"transactions"`
	Data           []byte                     `json:"data"`
	PrevBlockHash  []byte                     `json:"previous_block"`
	MerkleRootHash []byte                     `json:"merkle_root_hash"`
	Hash           []byte                     `json:"hash"`
}

func (b *Block) SetHash() {
	timestamp := []byte(fmt.Sprintf("%d", b.Timestamp))
	headers := bytes.Join([][]byte{b.Data, b.PrevBlockHash, timestamp, b.MerkleRootHash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) SetMerkleRoot() {
	if len(b.Transactions) == 0 {
		return
	}

	merkleHash, err := transaction.MerkleHashTransactions(b.Transactions)
	if err != nil {
		log.Printf("failed to generate merkle hash: %+v", err)
		return
	}

	b.MerkleRootHash = merkleHash
}

func NewBlock(transactions []*transaction.Transaction, data, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		Data:          data,
		PrevBlockHash: preBlockHash,
	}
	block.SetMerkleRoot()
	block.SetHash()

	return block
}
