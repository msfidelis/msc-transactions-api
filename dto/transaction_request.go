package dto

type TransactionRequest struct {
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
