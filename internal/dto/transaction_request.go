package dto

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Timestamp                  *int64   `json:"timestamp"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil || tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil || tr.Value == nil || tr.Signature == nil || tr.Timestamp == nil {
		return false
	}
	return true
}
