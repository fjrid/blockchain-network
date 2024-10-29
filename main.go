package main

import (
	"github.com/fjrid/blockchain-network/node"
	"github.com/fjrid/blockchain-network/server"
)

func main() {
	server.InitializeNetwork(node.InitNode())
}
