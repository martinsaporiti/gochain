package blockchain

import (
	"testing"
	"time"
)

func TestBlockchain_CreateMinerTransaction(t *testing.T) {

	blk := NewBlockchain("Node 500", "THE BLOCKCHAIN", 1)

	tests := map[string]struct {
		input *Blockchain
		want  *Transaction
	}{
		"should return a miner transaction": {
			input: blk,
			want: &Transaction{
				senderBlockchainAddress:    "Node 500",
				recipientBlockchainAddress: "THE BLOCKCHAIN",
				value:                      1.0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tx := tc.input.CreateMinerTransaction()

			if tx.senderBlockchainAddress != tc.want.senderBlockchainAddress {
				t.Errorf("CreateMinerTransaction() = %v, want %v", tx.senderBlockchainAddress, tc.want.senderBlockchainAddress)
			}

			if tx.recipientBlockchainAddress != tc.want.recipientBlockchainAddress {
				t.Errorf("CreateMinerTransaction() = %v, want %v", tx.recipientBlockchainAddress, tc.want.recipientBlockchainAddress)
			}

			if tx.value != tc.want.value {
				t.Errorf("CreateMinerTransaction() = %v, want %v", tx.value, tc.want.value)
			}

			if tx.timestamp == 0 {
				t.Errorf("CreateMinerTransaction() = %v, want %v", tx.timestamp, tc.want.timestamp)
			}
		})
	}
}

func TestBlockchain_CreateBlock(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	value := float32(200)
	timestamp := int64(1654369662)

	tx := &Transaction{
		senderBlockchainAddress:    sba,
		recipientBlockchainAddress: rba,
		value:                      value,
		timestamp:                  timestamp,
	}

	txs := make([]*Transaction, 0)
	txs = append(txs, tx)

	type input struct {
		number       int64
		nonce        int
		previousHash [32]byte
		transactions []*Transaction
	}

	blockchain := NewBlockchain("Node 500", "THE BLOCKCHAIN", 1)

	tests := map[string]struct {
		input input
		want  *Block
	}{
		"should create a new block": {
			input: input{
				number:       1,
				nonce:        123,
				previousHash: blockchain.LastBlock().Hash(),
				transactions: txs,
			},

			want: &Block{
				number:       1,
				nonce:        123,
				previousHash: blockchain.LastBlock().Hash(),
				transactions: txs,
			},
		},

		"sould not create a block": {
			input: input{
				number:       1,
				nonce:        0,
				previousHash: blockchain.LastBlock().Hash(),
				transactions: txs,
			},
			want: nil,
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			blk := blockchain.CreateBlock(tc.input.number, tc.input.nonce, tc.input.previousHash, tc.input.transactions)

			if blk != nil {

				// This checks the blocks was added to the blockchain.
				if blockchain.LastBlock().Hash() != blk.Hash() {
					t.Errorf("CreateBlock() = %v, want %v", blockchain.LastBlock().Hash(), blk.Hash())
				}

				if blk.number != tc.want.number {
					t.Errorf("CreateBlock() = %v, want %v", blk.number, tc.want.number)
				}

				if blk.nonce != tc.want.nonce {
					t.Errorf("CreateBlock() = %v, want %v", blk.nonce, tc.want.nonce)
				}

				if blk.previousHash != tc.want.previousHash {
					t.Errorf("CreateBlock() = %v, want %v", blk.previousHash, tc.want.previousHash)
				}

				if len(blk.transactions) != len(tc.want.transactions) {
					t.Errorf("CreateBlock() = %v, want %v", len(blk.transactions), len(tc.want.transactions))
				}

				if blk.transactions[0].senderBlockchainAddress != tc.want.transactions[0].senderBlockchainAddress {
					t.Errorf("CreateBlock() = %v, want %v", blk.transactions[0].senderBlockchainAddress, tc.want.transactions[0].senderBlockchainAddress)
				}

				if blk.transactions[0].recipientBlockchainAddress != tc.want.transactions[0].recipientBlockchainAddress {
					t.Errorf("CreateBlock() = %v, want %v", blk.transactions[0].recipientBlockchainAddress, tc.want.transactions[0].recipientBlockchainAddress)
				}

				if blk.transactions[0].value != tc.want.transactions[0].value {
					t.Errorf("CreateBlock() = %v, want %v", blk.transactions[0].value, tc.want.transactions[0].value)
				}

				if blk.transactions[0].timestamp != tc.want.transactions[0].timestamp {
					t.Errorf("CreateBlock() = %v, want %v", blk.transactions[0].timestamp, tc.want.transactions[0].timestamp)
				}
			}

			if blk == nil && tc.want != nil {
				t.Errorf("CreateBlock() = %v, want %v", blk, tc.want)
			}
		})

	}
}

func TestBlockchain_AddProposedBlockFromNetwork(t *testing.T) {

	blockchain := NewBlockchain("Node 500", "THE BLOCKCHAIN", 1)

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	value := float32(200)
	timestamp := int64(1654369662)

	tx := &Transaction{
		senderBlockchainAddress:    sba,
		recipientBlockchainAddress: rba,
		value:                      value,
		timestamp:                  timestamp,
	}

	txs := make([]*Transaction, 0)
	txs = append(txs, tx)

	ok := false
	nonce := 0
	for !ok {
		nonce++
		ok = blockchain.validProof(2, nonce, blockchain.LastBlock().Hash(), txs, 1)
	}

	blk1 := &Block{
		number:       2,
		nonce:        nonce,
		previousHash: blockchain.LastBlock().Hash(),
		transactions: txs,
		timestamp:    time.Now().Unix(),
	}

	blk2 := &Block{
		number:       2,
		nonce:        nonce,
		previousHash: blockchain.LastBlock().Hash(),
		transactions: txs,
		timestamp:    time.Now().Add(time.Second * time.Duration(-5)).Unix(),
	}

	tests := map[string]struct {
		input *Block
		want  bool
	}{
		"should add a block": {
			input: blk1,
			want:  true,
		},

		"should add the block with same number": {
			input: blk2,
			want:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			added := blockchain.AddProposedBlockFromNetwork(tc.input)
			if added != tc.want {
				t.Errorf("AddProposedBlockFromNetwork() = %v, want %v", added, tc.want)
			}
		})
	}
}
