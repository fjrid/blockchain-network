package main

import (
	"context"

	"github.com/fjrid/blockchain-network/network"
)

func main() {
	network.InitializeNetwork(context.Background())
}
