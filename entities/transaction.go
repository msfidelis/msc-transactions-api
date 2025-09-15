package entities

import "github.com/uptrace/bun"

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`
	IDClient      string `json:"id_client"`
	Amount        int64  `json:"amount"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	Date          string `json:"date"`
}
