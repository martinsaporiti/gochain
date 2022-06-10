package blockchain

import (
	"crypto/ecdsa"
	"reflect"
	"testing"

	"github.com/martinsaporiti/blockchain-sample/internal/blkcrypto"
	"github.com/martinsaporiti/blockchain-sample/internal/dto"
)

func TestTransactionPool_VerifyTransactionSignature(t *testing.T) {

	type input struct {
		senderPublicKey *ecdsa.PublicKey
		signature       *blkcrypto.Signature
		tx              *Transaction
	}

	tests := map[string]struct {
		input  input
		txPool *TransactionPool
		want   bool
	}{
		"should return true when the signature is valid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 200, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   true,
		},
		"should return false when the signature is invalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006c"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 200, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
		"should return false when the public key is inavalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce2"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 200, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
		"should return false when the sender is invalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrc", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 200, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
		"should return false when the receiver is invalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXS", 200, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
		"should return false when the amount is invalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 201, 1654369662),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
		"should return false when the timestamp is invalid": {
			input: input{
				senderPublicKey: blkcrypto.PublicKeyFromString("aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"),
				signature:       blkcrypto.SignatureFromString("1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"),
				tx:              NewTransaction("15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk", "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW", 200, 1654369661),
			},
			txPool: NewTransactionPool(nil),
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.txPool.verifyTransactionSignature(tc.input.senderPublicKey, tc.input.signature, tc.input.tx)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTransactionPool_AddAndVerifyTransaction(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	invalid_sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrl"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	spk := "aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"
	value := float32(200)
	timestamp := int64(1654369662)
	sig := "1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"

	tests := map[string]struct {
		input          *dto.TransactionRequest
		want           bool
		expectedLenght int
	}{
		"should return true when the transaction is valid": {
			input: &dto.TransactionRequest{
				SenderBlockchainAddress:    &sba,
				RecipientBlockchainAddress: &rba,
				SenderPublicKey:            &spk,
				Value:                      &value,
				Timestamp:                  &timestamp,
				Signature:                  &sig,
			},
			want:           true,
			expectedLenght: 1,
		},
		"should return false when the sendeer blockchain address is not valid in the transaction": {
			input: &dto.TransactionRequest{
				SenderBlockchainAddress:    &invalid_sba,
				RecipientBlockchainAddress: &rba,
				SenderPublicKey:            &spk,
				Value:                      &value,
				Timestamp:                  &timestamp,
				Signature:                  &sig,
			},
			want:           false,
			expectedLenght: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			txPool := NewTransactionPool(nil)
			got := txPool.AddAndVerifyTransaction(tc.input)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
			if len(txPool.transactions) != tc.expectedLenght {
				t.Errorf("got %v, want %v", len(txPool.transactions), tc.expectedLenght)
			}
		})
	}

}

func TestTransactionPool_Copy(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	spk := "aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"
	value := float32(200)
	timestamp := int64(1654369662)
	sig := "1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"

	tests := map[string]struct {
		input func() *TransactionPool
		want  []*Transaction
	}{
		"should return a copy of the transaction pool": {
			input: func() *TransactionPool {
				txPool := NewTransactionPool(nil)
				txPool.AddAndVerifyTransaction(&dto.TransactionRequest{
					SenderBlockchainAddress:    &sba,
					RecipientBlockchainAddress: &rba,
					SenderPublicKey:            &spk,
					Signature:                  &sig,
					Value:                      &value,
					Timestamp:                  &timestamp,
				})
				return txPool
			},
			want: []*Transaction{
				{
					senderBlockchainAddress:    sba,
					recipientBlockchainAddress: rba,
					value:                      value,
					timestamp:                  timestamp,
				},
			},
		},
		"should return an empty slice when the transaction pool is empty": {
			input: func() *TransactionPool {
				return NewTransactionPool(nil)
			},
			want: []*Transaction{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			txPool := tc.input()
			got := txPool.Copy()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
			if txPool.Length() > 0 {
				t.Errorf("got %v, want %v", txPool.Length(), 0)
			}
		})
	}
}

func TestTransactionPool_UpdateFromBlock(t *testing.T) {

	sba := "15TZoyyxFmeTXJGjYwX1X3ARtXX94BbFrk"
	rba := "1CHD4Jjqsak4RV5JHAdYZ9CKY1dQe4tkXW"
	spk := "aed86cb86fe477183b5a9f452d2a7d26e81c9ce16b123c699cbc4cc61cb7111df24ceebf3cb316b5740c19451ed390e9a9b1f5070cef639808af535570b01ce1"
	value := float32(200)
	timestamp := int64(1654369662)
	sig := "1dda674ee41218569a870993fe32fd3bf6a7bda3c657d4fae8fd0898adfdaf59eced8a5dd315f239c056197cf2b9a7a0bc4e941853be66e7e8a8427b14b9006b"

	sba1 := "14Zd3uAwyovSpuscsE7oYhPUko346z6Zne"
	rba1 := "1JkfWtkFzLHKoa33Vimaxcctc3z2HNWoet"
	spk1 := "b5238d60aa4f89da3b2fd6e1f61912ed567b9c92f7b920ef15c2677c9638b497669820ff0f419d9104c82bce464af6c0467498f092bfcd938dc5060c9e83729b"
	value1 := float32(450)
	timestamp1 := int64(1654689626)
	sig1 := "cc08d0df645f1a87a97f1974bd86d3cc1d647849d961b1a5675bd3d56bb8e79e99a17255f1f7a4e09900d84837bf906c5669ad0953b56dd7d6f712328070cb97"

	tests := map[string]struct {
		input  func() *Block
		txPool func() *TransactionPool
		want   []*Transaction
	}{
		"should return a slice with only one transaction ": {
			input: func() *Block {
				return &Block{
					transactions: []*Transaction{
						{
							senderBlockchainAddress:    sba,
							recipientBlockchainAddress: rba,
							value:                      value,
							timestamp:                  timestamp,
						},
					},
				}
			},
			txPool: func() *TransactionPool {
				txPool := NewTransactionPool(nil)
				txPool.AddAndVerifyTransaction(&dto.TransactionRequest{
					SenderBlockchainAddress:    &sba,
					RecipientBlockchainAddress: &rba,
					SenderPublicKey:            &spk,
					Signature:                  &sig,
					Value:                      &value,
					Timestamp:                  &timestamp,
				})

				txPool.AddAndVerifyTransaction(&dto.TransactionRequest{
					SenderBlockchainAddress:    &sba1,
					RecipientBlockchainAddress: &rba1,
					SenderPublicKey:            &spk1,
					Signature:                  &sig1,
					Value:                      &value1,
					Timestamp:                  &timestamp1,
				})
				return txPool
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			block := tc.input()
			txPool := tc.txPool()
			txPool.UpdateFromBlock(block)

			if len(txPool.transactions) != 1 {
				t.Errorf("got %v, want %v", len(txPool.transactions), 1)
			}

			if txPool.Transactions()[0].senderBlockchainAddress != sba1 {
				t.Errorf("got %v, want %v", txPool.Transactions()[0].senderBlockchainAddress, sba1)
			}

		})
	}
}
