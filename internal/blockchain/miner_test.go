package blockchain

import (
	"sync"
	"testing"
	"time"
)

func TestMiner_mineBlockComplete(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	value := float32(200)
	timestamp := int64(1654369662)
	tx := NewTransaction(sba, rba, value, timestamp)

	blockchain := NewBlockchain("a node name", "a node address", 1)
	txPool := NewTransactionPool(nil)
	txPool.transactions = make(map[string]*Transaction)
	txPool.transactions[tx.ID()] = tx

	startMining := make(chan bool)
	newBlockMined := make(chan *Block)

	miner := NewMiner(blockchain, txPool, startMining, newBlockMined)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer close(newBlockMined)
		newBlock := <-newBlockMined
		if newBlock == nil {
			t.Error("newBlock is nil")
		}

	}(&wg)

	go func() {
		miner.SignalStartMining()
	}()

	startMining <- true
	wg.Wait()
}

func TestMiner_mineBlockCanceled(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	value := float32(200)
	timestamp := int64(1654369662)
	tx := NewTransaction(sba, rba, value, timestamp)

	blockchain := NewBlockchain("a node name", "a node address", 10)
	txPool := NewTransactionPool(nil)
	txPool.transactions = make(map[string]*Transaction)
	txPool.transactions[tx.ID()] = tx

	startMining := make(chan bool)
	newBlockMined := make(chan *Block)
	miner := NewMiner(blockchain, txPool, startMining, newBlockMined)

	go func() {
		startMining <- true
	}()

	go func() {
		time.Sleep(time.Second * 3)
		miner.SignalCancelMining()
		close(startMining)
	}()

}
