package blockchain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
	timestamp                  int64
}

func NewTransaction(sender string, recipient string, value float32, timestamp int64) *Transaction {
	return &Transaction{sender, recipient, value, timestamp}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	fmt.Printf("senderBlockchainAddress: %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipientBlockchainAddress: %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value: %.1f\n", t.value)
	fmt.Printf("timestamp: %d\n", t.timestamp)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
		Timestamp int64   `json:"timestamp"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
		Timestamp: t.timestamp,
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		Sender    *string  `json:"sender_blockchain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
		Timestamp *int64   `json:"timestamp"`
	}{
		Sender:    &t.senderBlockchainAddress,
		Recipient: &t.recipientBlockchainAddress,
		Value:     &t.value,
		Timestamp: &t.timestamp,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) ID() string {
	return fmt.Sprintf("%x_%x_%x", t.senderBlockchainAddress, t.recipientBlockchainAddress, t.timestamp)
}

func (t *Transaction) Equal(tx *Transaction) bool {
	return t.ID() == tx.ID()
}
