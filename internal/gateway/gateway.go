package gateway

import (
	"sync"

	"github.com/martinsaporiti/blockchain-sample/internal/blockchain"
)

type Gateway interface {
	StartSyncNeighbors()
	NotifyNeighbors(endpoint, method string, message interface{})
	GetChains(wg *sync.WaitGroup) chan []*blockchain.Block
	NumberOfNeighbords() int
}
