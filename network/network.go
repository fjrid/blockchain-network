package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/node"
)

type network struct {
	node *node.Node
}

func (n *network) HandleAddBlock(w http.ResponseWriter, r *http.Request) {
	body := &AddBlockRequest{}
	json.NewDecoder(r.Body).Decode(body)

	n.node.AddBlock(body.Data)
	fmt.Fprintf(w, "Block added")
}

func (n *network) HandleReceiveNewBlock(w http.ResponseWriter, r *http.Request) {
	block := &block.Block{}
	json.NewDecoder(r.Body).Decode(block)

	err := n.node.ReceiveBlock(block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Block added")
}

func (n *network) HandleGetBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.node.GetBlockchains())
}

func (n *network) HandleAddPeer(w http.ResponseWriter, r *http.Request) {
	body := &AddNewPeerRequest{}
	json.NewDecoder(r.Body).Decode(&body)

	if body.Address == "" {
		http.Error(w, "address is empty", http.StatusBadRequest)
		return
	}

	n.node.AddNewPeer(body.Address)
}

func (n *network) HandleConnectToPeer(w http.ResponseWriter, r *http.Request) {
	body := &AddNewPeerRequest{}
	json.NewDecoder(r.Body).Decode(&body)

	if body.Address == "" {
		http.Error(w, "address is empty", http.StatusBadRequest)
		return
	}

	err := n.node.ConnectToPeer(body.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Success to add peer")
}

func (n *network) HandleGetPeers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(n.node.GetPeers())
}

func InitializeNetwork(node *node.Node, port string) {
	h := network{
		node: node,
	}

	http.HandleFunc("POST /add-block", h.HandleAddBlock)
	http.HandleFunc("POST /receive-block", h.HandleReceiveNewBlock)
	http.HandleFunc("POST /add-peer", h.HandleAddPeer)
	http.HandleFunc("POST /connect-to-peer", h.HandleConnectToPeer)
	http.HandleFunc("GET /blocks", h.HandleGetBlocks)
	http.HandleFunc("GET /peers", h.HandleGetPeers)

	log.Printf("App started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start node: %+v", err)
	}
}
