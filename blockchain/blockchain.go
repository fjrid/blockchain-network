package blockchain

import (
	"bytes"
	"errors"

	"github.com/fjrid/blockchain-network/block"
)

type Blockchain struct {
	blocks []*block.Block
}

func (bc *Blockchain) AddBlock(data string) *block.Block {
	block := block.NewBlock([]byte(data), bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, block)
	return block
}

func (bc *Blockchain) ReceiveBlock(block *block.Block) error {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	if !bytes.Equal(prevBlock.Hash, block.PrevBlockHash) && len(bc.blocks) > 1 {
		return errors.New("invalid block")
	}

	bc.blocks = append(bc.blocks, block)
	return nil
}

func (bc *Blockchain) GetBlocks() []*block.Block {
	return bc.blocks
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*block.Block{block.GenesisBlock()}}
}
