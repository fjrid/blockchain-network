package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fjrid/blockchain-network/network"
	"github.com/fjrid/blockchain-network/node"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("need to define port on first arg")
	}

	port := fmt.Sprintf(":%s", os.Args[1])
	network.InitializeNetwork(node.InitNode(port), port)
}
