package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var onceDB sync.Once
var onceDBPGX sync.Once
var onceBun sync.Once
var pgxInstance *sql.DB
var dbInstance *sql.DB
var BunInstance *bun.DB

func GetDBConn() *sql.DB {
	onceDB.Do(func() {
		var err error
		connectionString := getDBUrl()
		dbInstance, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("Error in database connection: %v", err)
		}

		// Verifica a conexão
		err = dbInstance.Ping()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	})
	return dbInstance
}

// func GetPGXPool() {
// 	onceDBPGX.Do(func() {
// 		var err error
// 		config, err := pgxpool.ParseConfig(getDBUrl())
// 		if err != nil {
// 			panic(err)
// 		}
// 		pgxInstance = stdlib.OpenDB(*config)
// 	})

// }

// Retorna a conexão com o database em utilizando uma estratégia de Singleton
func GetPGX() *sql.DB {
	onceDBPGX.Do(func() {
		var err error
		config, err := pgx.ParseConfig(getDBUrl())
		if err != nil {
			log.Fatalf("Error parsing database config: %v", err)
		}

		// Configurações de timeout e pool de conexão
		config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
		config.ConnectTimeout = 30 * time.Second

		// Log da configuração para debug
		log.Printf("Connecting to PostgreSQL: %s:%d", config.Host, config.Port)

		pgxInstance = stdlib.OpenDB(*config)

		// Pool de conexões mais conservador
		pgxInstance.SetMaxOpenConns(5)
		pgxInstance.SetMaxIdleConns(2)
		pgxInstance.SetConnMaxLifetime(30 * time.Minute)
		pgxInstance.SetConnMaxIdleTime(5 * time.Minute)

		// Testa a conexão
		if err = pgxInstance.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		log.Printf("Successfully connected to PostgreSQL")
	})
	return pgxInstance
}

func getDBUrl() string {
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	schema := os.Getenv("DATABASE_DB")

	// Log dos valores para debug
	log.Printf("Database connection config - Host: %s, Port: %s, User: %s, DB: %s", host, port, user, schema)

	// String de conexão com parâmetros de timeout e configurações de rede
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=30&statement_timeout=30000&idle_in_transaction_session_timeout=30000",
		user, pass, host, port, schema,
	)

	return connectionString
}

func GetDB() *bun.DB {
	onceBun.Do(func() {
		conn := GetPGX()
		// conn := GetDBConn()
		BunInstance = bun.NewDB(conn, pgdialect.New(), bun.WithDiscardUnknownColumns())
	})
	return BunInstance
}
