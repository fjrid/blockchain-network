package transaction

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fjrid/blockchain-network/util"
)

type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func NewTransaction(from, to string, amount float64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (t *Transaction) Hash() []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%f", t.From, t.To, t.Amount)))
	return hash[:]
}

func (t *Transaction) RLP() ([]byte, error) {
	return rlp.EncodeToBytes(
		bytes.Join([][]byte{[]byte(t.From), []byte(t.To), util.Float64ToBytes(t.Amount)}, []byte{}),
	)
}
