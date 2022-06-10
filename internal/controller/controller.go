package controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/martinsaporiti/blockchain-sample/internal/blockchain"
	"github.com/martinsaporiti/blockchain-sample/internal/config"
	"github.com/martinsaporiti/blockchain-sample/internal/dto"
	"github.com/martinsaporiti/blockchain-sample/internal/gateway"
)

const (
	MINING_SENDER = "THE BLOCKCHAIN"
)

type Controller interface {
	GetBlockchain() *blockchain.Blockchain
	CreateTransaction(tx *dto.TransactionRequest) bool
	AddTransaction(tr *dto.TransactionRequest) bool
	GetTransactions() []*blockchain.Transaction
	AddProposedBlockFromNetwork(block *blockchain.Block)
	CalculateTotalAmount(blockchainAddress string) float32
}

type controller struct {
	blockchainAddress    string
	blockchain           *blockchain.Blockchain
	gateway              gateway.Gateway
	txPool               *blockchain.TransactionPool
	miner                blockchain.Miner
	nodeName             string
	startMiningChannel   chan bool
	newBlockMinedChannel chan *blockchain.Block
}

func New(config config.Config) Controller {
	gtw := gateway.New(config.Port)
	nodeName := MINING_SENDER + " " + strconv.FormatInt(int64(config.Port), 10)
	blchain := blockchain.NewBlockchain(nodeName, config.BlockchainAddress, config.MiningDifficulty)

	startMiningChannel := make(chan bool)
	newBlockMinedChannel := make(chan *blockchain.Block)
	txPool := blockchain.NewTransactionPool(startMiningChannel)

	miner := blockchain.NewMiner(blchain, txPool, startMiningChannel, newBlockMinedChannel)

	ctrl := &controller{
		blockchainAddress:    config.BlockchainAddress,
		blockchain:           blchain,
		gateway:              gtw,
		txPool:               txPool,
		nodeName:             nodeName,
		miner:                miner,
		startMiningChannel:   startMiningChannel,
		newBlockMinedChannel: newBlockMinedChannel,
	}

	ctrl.start()
	return ctrl
}

func (c *controller) start() {
	c.gateway.StartSyncNeighbors()
	c.updateBlockchainFromNetwork()
	go c.miner.SignalStartMining()
	go c.newBlockMined(c.newBlockMinedChannel)
}

// updateBlockchainFromNetwork - Updates the blockchain from the network.
func (c *controller) updateBlockchainFromNetwork() {
	maxLength := len(c.blockchain.Chain()) - 1
	var longestChain []*blockchain.Block = nil

	wg := sync.WaitGroup{}
	wg.Add(c.gateway.NumberOfNeighbords())

	chainsChan := c.gateway.GetChains(&wg)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(chainsChan)
	}(&wg)

	for chain := range chainsChan {
		log.Printf("The length of the chain received is %d and mine is: %d", len(chain), len(c.blockchain.Chain()))
		if len(chain) > maxLength && c.blockchain.IsValidChain(chain) {
			maxLength = len(chain)
			longestChain = chain
		}
	}

	if longestChain != nil {
		c.blockchain.SetChain(longestChain)
		// TODO: Evaluate this:
		c.txPool.UpdateFromBlock(c.blockchain.LastBlock())
		log.Printf("New chain is %d blocks long and is valid", len(longestChain))
	}
}

// GetBlockchain - Returns the blockchain.
// This method is called by the neighbors.
func (c *controller) GetBlockchain() *blockchain.Blockchain {
	return c.blockchain
}

// CreateTransaction - Creates a new transaction and adds it to the pool.
// Notifys the neighbors of the new transaction.
func (c *controller) CreateTransaction(tx *dto.TransactionRequest) bool {
	done := c.txPool.AddAndVerifyTransaction(tx)
	// if the transaction was added to the pool, we have to notify the neighbors,
	// broadcasting the transaction
	if done {
		c.gateway.NotifyNeighbors("transactions", http.MethodPut, tx)
	}

	return done
}

// AddTransaction - Adds a transaction to the pool.
// This method is called by the neighbors.
func (c *controller) AddTransaction(tr *dto.TransactionRequest) bool {
	return c.txPool.AddAndVerifyTransaction(tr)
}

// GetTransactions - Returns the transactions of the pool.
func (c *controller) GetTransactions() []*blockchain.Transaction {
	return c.txPool.Transactions()
}

// AddProposedBlockFromNetwork - Adds a block to the blockchain.
// This method is called by the neighbors.
// If the proposed block is valid, it is added to the blockchain and all the transactions in the block are removed
// from the pool.
func (c *controller) AddProposedBlockFromNetwork(block *blockchain.Block) {
	added := c.blockchain.AddProposedBlockFromNetwork(block)
	if added {
		// All the transactions in the block are removed from the pool.
		c.txPool.UpdateFromBlock(block)
		// Stops the miner (current mining operation).
		c.miner.SignalCancelMining()
	}
}

// CalculateTotalAmount - Returns the total amount of USD per a given address.
func (c *controller) CalculateTotalAmount(blockchainAddress string) float32 {
	return c.blockchain.CalculateTotalAmount(blockchainAddress)
}

// newBlockMined - Called when a new block is mined.
// Notifies the neighbors of the new block.
func (c *controller) newBlockMined(newBlockMinedChannel chan *blockchain.Block) {
	for block := range newBlockMinedChannel {
		c.gateway.NotifyNeighbors("add_block", http.MethodPost, block)
	}
}
