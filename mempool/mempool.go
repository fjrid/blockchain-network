package mempool

import (
	"sort"
	"sync"

	"github.com/fjrid/blockchain-network/transaction"
)

type (
	Mempool struct {
		sync.Mutex
		transactions []*transaction.Transaction
	}
)

func NewMempool() *Mempool {
	return &Mempool{
		transactions: make([]*transaction.Transaction, 0),
	}
}

func (m *Mempool) AddTransaction(tx *transaction.Transaction) {
	m.Lock()
	defer m.Unlock()

	m.transactions = append(m.transactions, tx)

	sort.Slice(m.transactions, func(i, j int) bool {
		return m.transactions[i].GasPrice > m.transactions[j].GasPrice
	})
}

func (m *Mempool) TakeTransaction(n int) []*transaction.Transaction {
	m.Lock()
	defer m.Unlock()

	if len(m.transactions) <= n {
		results := m.transactions
		m.transactions = make([]*transaction.Transaction, 0)
		return results
	}

	results := m.transactions[:n]
	m.transactions = m.transactions[n:]

	return results
}

func (m *Mempool) GetTransactions() []*transaction.Transaction {
	return m.transactions
}
