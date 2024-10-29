package blockchain

import "github.com/fjrid/blockchain-network/block"

type Blockchain struct {
	blocks []*block.Block
}

func (bc *Blockchain) AddBlock(data string) {
	bc.blocks = append(bc.blocks, block.NewBlock([]byte(data), bc.blocks[len(bc.blocks)-1].Hash))
}

func (bc *Blockchain) GetBlocks() []*block.Block {
	return bc.blocks
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*block.Block{block.GenesisBlock()}}
}
