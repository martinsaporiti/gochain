package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MINING_REWARD = 1.0
)

type Blockchain struct {
	blockchainAddress string
	difficulty        int
	chain             []*Block
	txPool            *TransactionPool
	nodeName          string
	mux               sync.Mutex
}

func NewBlockchain(nodeName string, blockchainAddress string, miningDificulty int) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.difficulty = miningDificulty
	bc.CreateBlock(1, 0, b.Hash(), nil)
	bc.nodeName = nodeName
	return bc
}

func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

func (bc *Blockchain) SetChain(chain []*Block) {
	bc.chain = chain
}

func (bc *Blockchain) CreateMinerTransaction() *Transaction {
	return NewTransaction(bc.nodeName, bc.blockchainAddress, MINING_REWARD, time.Now().Unix())
}

// CreateBlock creates a new block in the blockchain
// The nonce is a number that is used to verify the block
// The previousHash is the hash of the previous block
// The transactions is a list of transactions in the block
func (bc *Blockchain) CreateBlock(number int64, nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	// TODO: Check this condition:
	if nonce != 0 && len(transactions) == 0 {
		log.Println("No transactions, skipping block creation")
		return nil
	}

	log.Printf("Creating block: %d with %d transactions", nonce, len(transactions))
	b := NewBlock(number, nonce, previousHash, transactions)
	if bc.addBlock(b) {
		return b
	}

	return nil
}

// addBlock - Adds a new block to the blockchain.
// If the block is valid, it will be added to the blockchain.
// If the block is invalid, it will be ignored.
func (bc *Blockchain) addBlock(block *Block) bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.chain) == 0 {
		bc.chain = append(bc.chain, block)
		return true
	}

	if block.previousHash != bc.LastBlock().Hash() {
		return false
	}

	bc.chain = append(bc.chain, block)
	return true
}

func (bc *Blockchain) replaceLastBlock(block *Block) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.chain[len(bc.chain)-1] = block
}

// LastBlock - Returns the last block of the blockchain
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// validProof - Validates the Proof.
// Returns true if the Proof is valid, false otherwise.
func (bc *Blockchain) validProof(number int64, nonce int, previousHash [32]byte, transactions []*Transaction,
	difficulty int) bool {

	zeros := strings.Repeat("0", difficulty)
	guessBlock := &Block{
		number: number, nonce: nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros

}

// AddProposedBlockFromNetwork - Adds a new block from the network
// If the block is valid, it will be added to the blockchain
// If the block is invalid, it will be ignored
func (bc *Blockchain) AddProposedBlockFromNetwork(proposedBlock *Block) bool {
	currentLastBlock := bc.LastBlock()
	newBlockAccepted := false
	if proposedBlock.Number() == currentLastBlock.Number()+1 {
		log.Printf("Adding new block from network: %d", proposedBlock.Number())
		previousHash := currentLastBlock.Hash()
		if bc.validProof(proposedBlock.Number(), proposedBlock.Nonce(), previousHash, proposedBlock.Transactions(),
			bc.difficulty) {
			bc.addBlock(proposedBlock)
			log.Printf("New block added: %d", proposedBlock.Number())
			newBlockAccepted = true
		}
	}

	if proposedBlock.Number() == currentLastBlock.Number() {
		log.Printf("Proposed block has the same number than our last block: %d",
			proposedBlock.Number())

		if proposedBlock.timestamp < currentLastBlock.timestamp {
			log.Printf("Proposed block is older than the current last block: %d", proposedBlock.Number())
			previousHash := currentLastBlock.PreviousHash()
			if bc.validProof(proposedBlock.Number(), proposedBlock.Nonce(), previousHash, proposedBlock.Transactions(),
				bc.difficulty) {
				log.Printf("Adding new block from network after removing current lastblock  block: %d",
					proposedBlock.Number())
				bc.replaceLastBlock(proposedBlock)
				newBlockAccepted = true
			}
		}
	}

	if newBlockAccepted {
		log.Printf("Proposed Block was added: %d", proposedBlock.Number())
		return true
	}

	log.Printf("Proposed block was ignored: %d", proposedBlock.Number())
	return false
}

// isValidChain - Validates the chain.
// Returns true if the chain is valid, false otherwise.
func (bc *Blockchain) IsValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			return false
		}

		if !bc.validProof(b.Number(), b.Nonce(), b.PreviousHash(), b.Transactions(), bc.difficulty) {
			log.Printf("Nonce is invalid: %d", b.Nonce())
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

// CalculateTotalAmount - Calculates the total amount for a Blockchain Address
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if t.recipientBlockchainAddress == blockchainAddress {
				totalAmount += value
			}

			if t.senderBlockchainAddress == blockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) Transactions() []*Transaction {
	return bc.txPool.Transactions()
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) Print() {
	for i, b := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("-", 20), i, strings.Repeat("-", 20))
		b.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 50))
}
