package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (n *Node) AddBlock(data string) {
	n.mux.Lock()
	defer n.mux.Unlock()

	block := n.blockchain.AddBlock(data)
	n.broadcastBlock(block)
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

func (n *Node) AddNewPeer(address string) {
	n.mux.Lock()
	defer n.mux.Unlock()

	if n.isPeerExist(address) {
		return
	}

	n.peers = append(n.peers, &Peer{Address: address})
	log.Printf("connected to new peer: %s", address)
}

func (n *Node) ConnectToPeer(address string) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	if n.isPeerExist(address) {
		return nil
	}

	jsonData, _ := json.Marshal(&Peer{Address: n.host})
	if _, err := http.Post(fmt.Sprintf("%s/add-peer", address), "application/json", bytes.NewBuffer(jsonData)); err != nil {
		return err
	}

	n.peers = append(n.peers, &Peer{Address: address})
	log.Printf("connected to new peer: %s", address)

	return nil
}

func (n *Node) GetPeers() []*Peer {
	return n.peers
}

func (n *Node) broadcastBlock(block *block.Block) {
	jsonData, _ := json.Marshal(block)

	for _, peer := range n.peers {
		log.Println("Sending to peer %s", peer.Address)
		if _, err := http.Post(fmt.Sprintf("%s/receive-block", peer.Address), "application/json", bytes.NewBuffer(jsonData)); err != nil {
			log.Printf("failed to broadcast block (#%s) to %s", string(block.Hash), peer.Address)
		}
	}
}

func (n *Node) isPeerExist(address string) bool {
	for _, peer := range n.peers {
		if peer.Address == address {
			return true
		}
	}
	return false
}

func InitNode(port string) *Node {
	return &Node{
		host:       fmt.Sprintf("http://localhost%s", port), // Need to implement NAT-PMP and UPnP
		blockchain: blockchain.NewBlockChain(),
		peers:      make([]*Peer, 0),
	}
}
