package services

import (
	"context"
	"log"
	"main/entities"
	"main/pkg/database"
	"time"
)

func Statement(id_client string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.GetDB()
	err := db.NewSelect().
		Model(&transactions).
		Where("id_client = ?", id_client).
		Limit(10).
		Order("date DESC").
		Scan(ctx)
	if err != nil {
		log.Printf("[%s] Error %s: %v", "Statement", id_client, err)
		return nil, err
	}

	return transactions, nil

}
