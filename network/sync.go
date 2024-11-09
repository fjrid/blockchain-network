package network

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fjrid/blockchain-network/block"
	dnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (n *network) setupSyncStream(ctx context.Context) {
	n.nat.Host.SetStreamHandler(n.syncPID, n.incomingSyncHandler)

	if !n.nat.IsServerMode {
		n.syncData(ctx)
	}
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
