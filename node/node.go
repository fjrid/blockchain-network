package node

import (
	"sync"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/blockchain"
)

type Node struct {
	blockchain *blockchain.Blockchain
	mux        sync.Mutex
}

func (n *Node) AddBlock(data string) {
	n.mux.Lock()
	defer n.mux.Unlock()

	n.blockchain.AddBlock(data)
}

func (n *Node) GetBlockchains() []*block.Block {
	return n.blockchain.GetBlocks()
}

func InitNode() *Node {
	return &Node{
		blockchain: blockchain.NewBlockChain(),
	}
}
