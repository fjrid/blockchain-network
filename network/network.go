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
	dnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type network struct {
	node *node.Node
	nat  *NAT

	topic   *pubsub.Topic
	syncPID protocol.ID
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

		if msg.ReceivedFrom.String() != n.nat.Host.ID().String() {
			n.HandleReceiveNewBlock(msg.Data)
		}
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

func (n *network) incomingSyncHandler(stream dnet.Stream) {
	log.Println("Incoming request to sync data")
	if err := json.NewEncoder(stream).Encode(n.node.GetBlockchains()); err != nil {
		log.Printf("failed to sent sync blockchain response: %+v", err)
	}

	defer stream.Close()
}

func (n *network) syncData(ctx context.Context) {
	log.Println("Process to sync data")

	isSuccess := false
	for _, peer := range n.nat.Host.Network().Peers() {
		if err := n.requestSyncData(ctx, peer); err != nil {
			log.Printf("failed to request sync data: %+v", err)
			continue
		}

		isSuccess = true
		break
	}

	if !isSuccess {
		log.Fatalf("Failed to sync data, cannot connect to network")
	}

	log.Println("Success to sync data")
}

func (n *network) requestSyncData(ctx context.Context, peer peer.ID) error {
	stream, err := n.nat.Host.NewStream(ctx, peer, n.syncPID)
	if err != nil {
		return fmt.Errorf("failed to request sync data (%s): %+v", peer.ShortString(), err)
	}

	defer stream.Close()

	blocks := make([]*block.Block, 0)
	if err := json.NewDecoder(stream).Decode(&blocks); err != nil {
		return fmt.Errorf("failed to decoding bocks: %+v", err)
	}

	n.node.SetBlockchain(blocks)

	return nil
}

func (n *network) setupSyncStream(ctx context.Context) {
	n.nat.Host.SetStreamHandler(n.syncPID, n.incomingSyncHandler)

	if !n.nat.IsServerMode {
		n.syncData(ctx)
	}
}

func InitializeNetwork(ctx context.Context) {
	nat, err := initializeNAT(ctx)
	if err != nil {
		log.Fatalf("failed to initialize nat: %+v", err)
	}

	port := fmt.Sprintf(":%d", nat.Port)
	h := network{
		node: node.InitNode(port),
		nat:  nat,
	}

	h.setupSyncStream(ctx)
	h.initializeGossip()

	http.HandleFunc("POST /add-block", h.HandleAddBlock)
	http.HandleFunc("GET /blocks", h.HandleGetBlocks)
	http.HandleFunc("GET /peers", h.HandleGetPeers)

	log.Printf("App started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start node: %+v", err)
	}
}
