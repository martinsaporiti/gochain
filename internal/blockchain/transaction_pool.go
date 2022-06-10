package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/martinsaporiti/blockchain-sample/internal/blkcrypto"
	"github.com/martinsaporiti/blockchain-sample/internal/dto"
)

type TransactionPool struct {
	transactions       map[string]*Transaction
	mux                sync.Mutex
	startMiningChannel chan bool
}

func NewTransactionPool(startMiningChannel chan bool) *TransactionPool {

	return &TransactionPool{
		transactions:       make(map[string]*Transaction),
		startMiningChannel: startMiningChannel,
	}

}

func (tp *TransactionPool) isNodeAddress(address string) bool {
	return strings.Contains("THE BLOCKCHAIN", address)
}

// AddAndVerifyTransaction - Adds a transaction to the transaction pool
// and verifies the signature of the transaction.
// Sends a message to the mining process to start mining.
func (tp *TransactionPool) AddAndVerifyTransaction(tr *dto.TransactionRequest) bool {
	sender := *tr.SenderBlockchainAddress
	recipient := *tr.RecipientBlockchainAddress
	value := *tr.Value
	timestamp := *tr.Timestamp
	senderPublicKey := blkcrypto.PublicKeyFromString(*tr.SenderPublicKey)
	signature := blkcrypto.SignatureFromString(*tr.Signature)

	t := NewTransaction(sender, recipient, value, timestamp)
	if tp.isNodeAddress(sender) {
		tp.Add(t)
		return true
	}

	if tp.verifyTransactionSignature(senderPublicKey, signature, t) {
		tp.Add(t)
		log.Println("action = add transaction, status = success")

		if len(tp.transactions) == 1 {
			go func() {
				tp.startMiningChannel <- true
			}()
		}
		return true
	}

	log.Println("action = add transaction, status = failed")
	log.Println("ERROR: Invalid transaction signature")
	return false
}

// verifyTransactionSignature - Verifies the signature of a transaction
func (tp *TransactionPool) verifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *blkcrypto.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

// Transactions - Returns a transaction from the transaction pool
func (tp *TransactionPool) Transactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	tp.mux.Lock()
	defer tp.mux.Unlock()
	for _, t := range tp.transactions {
		transactions = append(transactions, t)
	}
	return transactions
}

// Add - Adds a transaction to the transaction pool
func (tp *TransactionPool) Add(t *Transaction) {
	tp.mux.Lock()
	defer tp.mux.Unlock()
	tp.transactions[t.ID()] = t
}

// Copy - Returns a copy of the transaction pool
// Removes all transactions from the pool.
func (tp *TransactionPool) Copy() []*Transaction {
	tp.mux.Lock()
	defer tp.mux.Unlock()
	transactions := make([]*Transaction, 0)
	for _, t := range tp.transactions {
		transactions = append(transactions, t)
	}
	tp.transactions = make(map[string]*Transaction)
	return transactions
}

// UpdateFromBlock - Updates the transaction pool removing the transactions from a block
func (tp *TransactionPool) UpdateFromBlock(b *Block) {
	for _, t := range b.Transactions() {
		tp.remove(t)
	}
}

func (tp *TransactionPool) Length() int {
	return len(tp.transactions)
}

// Remove - Removes a transaction from the transaction pool
func (tp *TransactionPool) remove(t *Transaction) {
	tp.mux.Lock()
	defer tp.mux.Unlock()
	delete(tp.transactions, t.ID())
}
