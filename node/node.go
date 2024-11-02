package node

import (
	"fmt"
	"sync"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/blockchain"
)

type Peer struct {
	Address string `json:"address"`
}

type Node struct {
	host       string
	blockchain *blockchain.Blockchain
	peers      []*Peer
	mux        sync.Mutex
}

func (n *Node) AddBlock(data string) *block.Block {
	n.mux.Lock()
	defer n.mux.Unlock()

	block := n.blockchain.AddBlock(data)
	return block
}

func (n *Node) ReceiveBlock(block *block.Block) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	if err := n.blockchain.ReceiveBlock(block); err != nil {
		return err
	}
	return nil
}

func (n *Node) GetBlockchains() []*block.Block {
	return n.blockchain.GetBlocks()
}

func (n *Node) GetPeers() []*Peer {
	return n.peers
}

func (n *Node) SetBlockchain(blocks []*block.Block) {
	n.mux.Lock()
	defer n.mux.Unlock()

	n.blockchain.SetBlock(blocks)
}

func InitNode(port string) *Node {
	return &Node{
		host:       fmt.Sprintf("http://localhost%s", port), // Need to implement NAT-PMP and UPnP
		blockchain: blockchain.NewBlockChain(),
		peers:      make([]*Peer, 0),
	}
}
