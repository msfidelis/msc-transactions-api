package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	redis "github.com/redis/go-redis/v9"
)

var onceRedis sync.Once
var redisInstance *redis.Client

// GetClient retorna a instância do cliente Redis usando padrão Singleton
func GetClient() *redis.Client {
	onceRedis.Do(func() {
		host := os.Getenv("CACHE_HOST")
		if host == "" {
			host = "localhost" // valor padrão
		}

		portStr := os.Getenv("CACHE_PORT")
		if portStr == "" {
			portStr = "6379" // valor padrão
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid CACHE_PORT value: %v", err)
		}

		// Configuração do cliente Redis
		redisInstance = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: "", // sem senha por padrão
			DB:       0,  // banco de dados padrão
		})

		// Teste de conexão
		ctx := context.Background()
		_, err = redisInstance.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Printf("Connected to Redis at %s:%d", host, port)
	})
	return redisInstance
}

// Close fecha a conexão com o Redis
func Close() error {
	if redisInstance != nil {
		return redisInstance.Close()
	}
	return nil
}
