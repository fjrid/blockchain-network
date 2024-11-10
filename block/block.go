package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fjrid/blockchain-network/db"
	"github.com/fjrid/blockchain-network/mpt"
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

func (b *Block) SetMerkleRoot(db *db.DB) error {
	if len(b.Transactions) == 0 {
		return nil
	}

	merklePatriciaTrie := mpt.NewMerklePatriciaTrie(db)
	for i, tx := range b.Transactions {
		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return err
		}

		val, err := tx.RLP()
		if err != nil {
			return err
		}

		merklePatriciaTrie.Insert(key, val)
	}

	merklePatriciaTrie.Store()
	b.MerkleRootHash = merklePatriciaTrie.Root.Hash()

	log.Printf("new block added with merkle root: %s", hex.EncodeToString(b.MerkleRootHash))
	return nil
}

func NewBlock(db *db.DB, transactions []*transaction.Transaction, data, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		Data:          data,
		PrevBlockHash: preBlockHash,
	}
	block.SetMerkleRoot(db)
	block.SetHash()

	return block
}
