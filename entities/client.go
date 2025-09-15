package entities

import "github.com/uptrace/bun"

type Client struct {
	bun.BaseModel `bun:"table:clients,alias:u"`
	ID            string `json:"id" bun:"id_client,pk,autoincrement"`
	Balance       int64  `json:"balance"`
	Limit         int64  `json:"limits" bun:"limits"`
}
