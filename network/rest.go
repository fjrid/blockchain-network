package network

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fjrid/blockchain-network/transaction"
)

func (h *network) setRoute() {
	http.HandleFunc("POST /block", h.HandleAddBlock)
	http.HandleFunc("POST /transaction", h.HandleAddTransaction)
	http.HandleFunc("GET /blocks", h.HandleGetBlocks)
	http.HandleFunc("GET /transactions/pending", h.HandleGetPendingTransactions)
	http.HandleFunc("POST /mpt/node/check", h.HandleCheckMPTNode)
	http.HandleFunc("GET /peers", h.HandleGetPeers)
}

func (n *network) HandleAddBlock(w http.ResponseWriter, r *http.Request) {
	body := &AddBlockRequest{}
	json.NewDecoder(r.Body).Decode(body)

	block := n.node.AddBlock(body.Data)
	jsonByte, _ := json.Marshal(block)

	err := n.topic.Publish(context.Background(), jsonByte)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to broadcast message: %+v", err), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Block added")
}

func (n *network) HandleAddTransaction(w http.ResponseWriter, r *http.Request) {
	body := &transaction.Transaction{}
	json.NewDecoder(r.Body).Decode(body)

	txHash := n.node.AddTransaction(body)

	fmt.Fprintf(w, "Success add transaction, your transaction still pending to process (transaction hash: #%s)", hex.EncodeToString(txHash))
}

func (n *network) HandleGetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.node.GetPendingTransactions())
}

func (n *network) HandleGetBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.node.GetBlockchains())
}

func (n *network) HandleGetPeers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(n.node.GetPeers())
}

func (n *network) HandleCheckMPTNode(w http.ResponseWriter, r *http.Request) {
	body := &struct {
		Hash string `json:"hash"`
	}{}
	json.NewDecoder(r.Body).Decode(body)

	hashKey, err := hex.DecodeString(body.Hash)
	if err != nil {
		fmt.Fprintf(w, "failed to decode hash: %+v", err)
	}

	result, err := n.node.CheckMPTNode(hashKey)
	if err != nil {
		fmt.Fprintf(w, "failed to check MPT Node: %+v", err)
	}

	fmt.Fprintf(w, "success to check MPT Node: %v", result)
}
