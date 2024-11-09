package transaction

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fjrid/blockchain-network/util"
)

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   uint64 `json:"amount"`
	GasPrice uint64 `json:"gas_price"`
}

func NewTransaction(from, to string, amount, gasPrice uint64) *Transaction {
	return &Transaction{
		From:     from,
		To:       to,
		Amount:   amount,
		GasPrice: gasPrice,
	}
}

func (t *Transaction) Hash() []byte {
	rlp, _ := t.RLP()
	return crypto.Keccak256(rlp)
}

func (t *Transaction) RLP() ([]byte, error) {
	return rlp.EncodeToBytes(
		bytes.Join([][]byte{
			[]byte(t.From),
			[]byte(t.To),
			util.Uint64ToBytes(t.Amount),
			util.Uint64ToBytes(t.GasPrice),
		}, []byte{}),
	)
}
