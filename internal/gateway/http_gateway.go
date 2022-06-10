package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/martinsaporiti/blockchain-sample/internal/blockchain"
	"github.com/martinsaporiti/blockchain-sample/internal/network"
)

const (
	// TODO: Change to config variables
	BLOCKCHAIN_PORT_RANGE_START      = 5000
	BLOCKCHAIN_PORT_RANGE_END        = 5003
	NEIGHBOR_IP_RANGE_START          = 1
	NEIGHBOR_IP_RANGE_END            = 3
	BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC = 10
)

type httpGateway struct {
	port         uint16
	neighbors    []string
	muxNeighbors sync.Mutex
}

func New(port uint16) Gateway {
	return &httpGateway{
		port: port,
	}
}

func (g *httpGateway) NumberOfNeighbords() int {
	return len(g.neighbors)
}

func (g *httpGateway) syncNeighbors() {
	g.muxNeighbors.Lock()
	defer g.muxNeighbors.Unlock()
	g.setNeighbors()
}

func (g *httpGateway) StartSyncNeighbors() {
	g.syncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC, g.StartSyncNeighbors)
}

func (g *httpGateway) setNeighbors() {
	g.neighbors = network.FindNeighbors(
		network.GetHost(), g.port, NEIGHBOR_IP_RANGE_START,
		NEIGHBOR_IP_RANGE_END, BLOCKCHAIN_PORT_RANGE_START,
		BLOCKCHAIN_PORT_RANGE_END)

}

// NotifyNeighbors - Notify all neighbors about a new transaction or block
func (g *httpGateway) NotifyNeighbors(endpoint, method string, message interface{}) {
	log.Printf("Neighbors: %v\n", g.neighbors)
	for _, n := range g.neighbors {
		go func(n string) {
			finalEndpoint := fmt.Sprintf("http://%s/%s", n, endpoint)
			client := &http.Client{}
			var req *http.Request
			if message != nil {
				m, _ := json.Marshal(message)
				body := bytes.NewBuffer(m)
				req, _ = http.NewRequest(method, finalEndpoint, body)
			} else {
				req, _ = http.NewRequest(method, finalEndpoint, nil)
			}
			_, err := client.Do(req)
			if err != nil {
				log.Println("ERROR: Failed notifying neighbors")
			}
		}(n)
	}
}

// GetChains - Returns the chains of all neighbors
func (g *httpGateway) GetChains(wg *sync.WaitGroup) chan []*blockchain.Block {
	// this channel will be used to return the chains of all neighbors:
	chainsChann := make(chan []*blockchain.Block)

	// For each neighbor, get the chain and add it to the channel
	for _, n := range g.neighbors {
		go func(wg *sync.WaitGroup, neighbor string, chainsChann chan []*blockchain.Block) {
			defer wg.Done()
			endpoint := fmt.Sprintf("http://%s/chain", n)
			log.Printf("Calling to resolve conflics: endpoint %s\n", endpoint)
			resp, _ := http.Get(endpoint)
			if resp.StatusCode == 200 {
				var bcResp blockchain.Blockchain
				decoder := json.NewDecoder(resp.Body)
				if err := decoder.Decode(&bcResp); err != nil {
					log.Printf("Error decoding response %v", err)
				}
				chain := bcResp.Chain()
				chainsChann <- chain
			}
		}(wg, n, chainsChann)
	}

	return chainsChann
}
