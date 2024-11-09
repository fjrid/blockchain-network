package node

import (
	"fmt"
	"sync"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/blockchain"
	"github.com/fjrid/blockchain-network/mempool"
	"github.com/fjrid/blockchain-network/transaction"
)

type Peer struct {
	Address string `json:"address"`
}

type Node struct {
	host            string
	blockchain      *blockchain.Blockchain
	peers           []*Peer
	transactionPool *mempool.Mempool
	mux             sync.Mutex

	maximumTransaction int
}

func (n *Node) AddBlock(data string) *block.Block {
	n.mux.Lock()
	defer n.mux.Unlock()

	transactions := n.transactionPool.TakeTransaction(n.maximumTransaction)

	block := n.blockchain.AddBlock(transactions, data)
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

func (n *Node) AddTransaction(tx *transaction.Transaction) []byte {
	n.transactionPool.AddTransaction(tx)

	return tx.Hash()
}

func (n *Node) GetPendingTransactions() []*transaction.Transaction {
	return n.transactionPool.GetTransactions()
}

func InitNode(port string) *Node {
	return &Node{
		host:               fmt.Sprintf("http://localhost%s", port), // Need to implement NAT-PMP and UPnP
		blockchain:         blockchain.NewBlockChain(),
		transactionPool:    mempool.NewMempool(),
		peers:              make([]*Peer, 0),
		maximumTransaction: 2,
	}
}
