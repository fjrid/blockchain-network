package network

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fjrid/blockchain-network/node"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type network struct {
	node *node.Node
	nat  *NAT

	topic   *pubsub.Topic
	syncPID protocol.ID
}

func InitializeNetwork(ctx context.Context) {
	nat, err := initializeNAT(ctx)
	if err != nil {
		log.Fatalf("failed to initialize nat: %+v", err)
	}

	port := fmt.Sprintf(":%d", nat.Port)
	h := network{
		node:    node.InitNode(port),
		nat:     nat,
		syncPID: protocol.ID("/sync/1.0.0"),
	}

	h.setupSyncStream(ctx)
	h.initializeGossip()
	h.setRoute()

	log.Printf("App started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start node: %+v", err)
	}
}
