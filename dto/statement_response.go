package dto

import "main/entities"

type StatementResponse struct {
	Balance          Balance                `json:"balance"`
	LastTransactions []entities.Transaction `json:"last_transactions"`
}

type Balance struct {
	Total         int64  `json:"total"`
	DateStatement string `json:"date_statement"`
	Limit         int64  `json:"limits"`
}
