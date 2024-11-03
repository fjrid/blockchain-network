package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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

	merklePatriciaTrie := transaction.NewMerklePatriciaTrie()
	for _, tx := range b.Transactions {
		val, err := json.Marshal(tx)
		if err == nil {
			merklePatriciaTrie.Insert(tx.Hash(), val)
		}
	}

	b.MerkleRootHash = []byte(merklePatriciaTrie.Root.Hash())
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
