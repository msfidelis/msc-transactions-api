package dto

type TransactionResponse struct {
	Limit   int64 `json:"limits"`
	Balance int64 `json:"balance"`
}
