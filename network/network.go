package network

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fjrid/blockchain-network/block"
	"github.com/fjrid/blockchain-network/node"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type network struct {
	node *node.Node
	nat  *NAT

	topic *pubsub.Topic
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

func (n *network) HandleReceiveNewBlock(data []byte) error {
	log.Println("Receive new block")

	block := &block.Block{}
	json.Unmarshal(data, block)

	err := n.node.ReceiveBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (n *network) HandleGetBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.node.GetBlockchains())
}

func (n *network) HandleGetPeers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(n.node.GetPeers())
}

func (n *network) listeningTopic(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Printf("failed to checking subscription: %+v", err)
		}

		n.HandleReceiveNewBlock(msg.Data)
	}
}

func (n *network) initializeGossip() error {
	ps, err := pubsub.NewGossipSub(context.Background(), n.nat.Host)
	if err != nil {
		return fmt.Errorf("failed to create gossip sub: %+v", err)
	}

	topic, err := ps.Join("/v1/BROADCAST/BLOCK")
	if err != nil {
		return fmt.Errorf("failed to join topic: %+v", err)
	}

	n.topic = topic

	sub, err := n.topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel: %+v", err)
	}

	go n.listeningTopic(context.Background(), sub)

	return nil
}

func InitializeNetwork() {
	nat, err := initializeNAT(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize nat: %+v", err)
	}

	port := fmt.Sprintf(":%d", nat.Port)
	h := network{
		node: node.InitNode(port),
		nat:  nat,
	}
	h.initializeGossip()

	http.HandleFunc("POST /add-block", h.HandleAddBlock)
	http.HandleFunc("GET /blocks", h.HandleGetBlocks)
	http.HandleFunc("GET /peers", h.HandleGetPeers)

	log.Printf("App started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start node: %+v", err)
	}
}
