package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
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

func (b *Block) SetMerkleRoot() error {
	if len(b.Transactions) == 0 {
		return nil
	}

	merklePatriciaTrie := mpt.NewMerklePatriciaTrie()
	for i, tx := range b.Transactions {
		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return err
		}

		merklePatriciaTrie.Insert(key, tx.Hash())
	}

	b.MerkleRootHash = merklePatriciaTrie.Root.Hash()
	return nil
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
