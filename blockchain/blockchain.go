package blockchain

import (
	"bytes"
	"errors"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/db"
	"github.com/fjrid/blockchain-network/transaction"
)

type Blockchain struct {
	blocks []*block.Block
}

func (bc *Blockchain) AddBlock(db *db.DB, transaction []*transaction.Transaction, data string) *block.Block {
	prevHash := make([]byte, 0)

	if len(bc.blocks) > 0 {
		prevHash = bc.blocks[len(bc.blocks)-1].Hash
	}

	block := block.NewBlock(db, transaction, []byte(data), prevHash)
	bc.blocks = append(bc.blocks, block)
	return block
}

func (bc *Blockchain) ReceiveBlock(block *block.Block) error {
	if len(bc.blocks) > 1 {
		prevBlock := bc.blocks[len(bc.blocks)-1]
		if !bytes.Equal(prevBlock.Hash, block.PrevBlockHash) && len(bc.blocks) > 1 {
			return errors.New("invalid block")
		}
	}

	bc.blocks = append(bc.blocks, block)
	return nil
}

func (bc *Blockchain) GetBlocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) SetBlock(blocks []*block.Block) {
	bc.blocks = blocks
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*block.Block{}}
}
