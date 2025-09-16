package services

import (
	"context"
	"fmt"
	"log"
	"main/pkg/database"
	"main/pkg/redis"
	"strconv"
	"time"
)

func GetBalance(id_client string) (int64, error) {
	functionName := "GetBalance"
	ctx := context.Background()

	// Cache key pattern: balance:id_client
	cacheKey := fmt.Sprintf("balance:%s", id_client)

	redisClient := redis.GetClient()

	cachedBalance, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - converte string para int64 e retorna
		log.Printf("[%s] Cache hit for client %s, balance: %s", functionName, id_client, cachedBalance)
		balance, parseErr := strconv.ParseInt(cachedBalance, 10, 64)
		if parseErr != nil {
			log.Printf("[%s] Error parsing cached balance for client %s: %v", functionName, id_client, parseErr)
		} else {
			log.Printf("[%s] Cache hit for client %s, balance: %d", functionName, id_client, balance)
			return balance, nil
		}
	}

	// Cache miss - busca no banco de dados
	log.Printf("[%s] Cache miss for client %s, fetching from database", functionName, id_client)

	db := database.GetDB()
	client, err := FindClient(ctx, db, id_client)
	if err != nil {
		log.Printf("[%s] Error finding client %s in database: %v", functionName, id_client, err)
		return 0, err
	}

	// Salva no cache para próximas consultas (TTL de 5 minutos)
	balanceStr := strconv.FormatInt(client.Balance, 10)
	err = redisClient.Set(ctx, cacheKey, balanceStr, 5*time.Minute).Err()
	if err != nil {
		log.Printf("[%s] Error caching balance for client %s: %v", functionName, id_client, err)
	} else {
		log.Printf("[%s] Cached balance for client %s: %d", functionName, id_client, client.Balance)
	}

	return client.Balance, nil
}
