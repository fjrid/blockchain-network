package network

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fjrid/blockchain-network/block"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

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
