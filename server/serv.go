package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fjrid/blockchain-network/node"
)

type handler struct {
	node *node.Node
}

func (h *handler) HandleAddBlock(w http.ResponseWriter, r *http.Request) {
	body := &AddBlockRequest{}
	json.NewDecoder(r.Body).Decode(body)

	h.node.AddBlock(body.Data)
	fmt.Fprintf(w, "Block added")
}

func (h *handler) HandleGetBlocks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.node.GetBlockchains())
}

func InitializeNetwork(node *node.Node) {
	h := handler{node}

	http.HandleFunc("POST /add-block", h.HandleAddBlock)
	http.HandleFunc("GET /blocks", h.HandleGetBlocks)

	log.Println("App started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start node: %+v", err)
	}
}
