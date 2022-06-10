package blockchain

import (
	"context"
	"fmt"
	"log"
)

type Miner interface {
	SignalStartMining()
	SignalCancelMining()
}

type miner struct {
	blockchain           *Blockchain
	txPool               *TransactionPool
	startMiningChannel   chan bool
	newBlockMinedChannel chan *Block
	ctx                  context.Context
	cancelFn             context.CancelFunc
}

func NewMiner(blockchain *Blockchain, txPool *TransactionPool, startMiningChannel chan bool, newBlockMinedChannel chan *Block) Miner {
	return &miner{
		blockchain:           blockchain,
		txPool:               txPool,
		startMiningChannel:   startMiningChannel,
		newBlockMinedChannel: newBlockMinedChannel,
	}
}

func (m *miner) SignalCancelMining() {
	m.cancelFn()
}

// SignalStartMining - start mining
// miner waits for a signal to start mining. That signal is sent by the tx pool.
func (m *miner) SignalStartMining() {
	log.Println(">>> Starting Miner")
	for {
		<-m.startMiningChannel
		if m.txPool.Length() > 0 {
			fmt.Println("Mining...")
			m.mineBlock()
		}
	}
}

// mineBlock - mine a new block
func (m *miner) mineBlock() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	m.ctx = ctx
	m.cancelFn = cancel

	log.Println(">>>> 1. action = mining, status = Starting")
	transactions := m.txPool.Copy()
	m.printTxs(transactions)
	transactions = append(transactions, m.blockchain.CreateMinerTransaction())
	previousHash := m.blockchain.LastBlock().Hash()
	lastNumber := m.blockchain.LastBlock().Number()

	nonce := m.proofOfWork(ctx, lastNumber+1, transactions)

	// if nonce == -1, then the miner was cancelled.
	// if nonce != -1, then the miner mined a block.
	if nonce != -1 {
		newBlock := m.blockchain.CreateBlock(lastNumber+1, nonce, previousHash, transactions)
		if newBlock == nil {
			log.Println(">>>> 4. action = mining, status = Failed, Block not created")
		} else {
			log.Println("<<<< 2. action = mining, status = Success")
			log.Printf("Sending new block over the network: %d", newBlock.number)
			m.newBlockMinedChannel <- newBlock
		}
	}

	// Checks if the are new transactions in the pool and puts to run again the miner...
	defer func() {
		m.SignalStartMining()
	}()
}

func (m *miner) proofOfWork(ctx context.Context, number int64, transactions []*Transaction) int {
	log.Printf(">>> Starting Proof of Work for %d transactions", len(transactions))
	previousHash := m.blockchain.LastBlock().Hash()
	nonce := -1
	done := false
	for !done {
		nonce += 1
		select {
		case <-ctx.Done():
			// Context cancelled, new block from network added.
			log.Println("<<<< 3. action = mining, status = Canceled")
			return -1
		default:
			done = m.blockchain.validProof(number, nonce, previousHash, transactions, m.blockchain.difficulty)
		}
	}

	return nonce
}

func (m *miner) printTxs(transactions []*Transaction) {
	for _, tx := range transactions {
		tx.Print()
	}
}
