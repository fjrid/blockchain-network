package network

import "github.com/fjrid/blockchain-network/transaction"

type AddBlockRequest struct {
	Data         string                     `json:"data"`
	Transactions []*transaction.Transaction `json:"transactions"`
}

type AddNewPeerRequest struct {
	Address string `json:"address"`
}
