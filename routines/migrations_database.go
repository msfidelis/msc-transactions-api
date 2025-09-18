package routines

import (
	"database/sql"
	"fmt"
	"main/pkg/database"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver do banco de dados
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func DatabaseMigration() {
	consumerName := "DatabaseMigration"

	// Aguarda mais tempo para o PostgreSQL estar totalmente pronto
	fmt.Printf("[%s] Aguardando PostgreSQL estar pronto...\n", consumerName)
	time.Sleep(30 * time.Second)

	// Retry logic para conexão com o banco
	var db *sql.DB
	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		db = database.GetPGX()
		if db != nil {
			// Testa a conexão
			if err = db.Ping(); err == nil {
				fmt.Printf("[%s] Conexão com PostgreSQL estabelecida na tentativa %d\n", consumerName, i+1)
				break
			}
			fmt.Printf("[%s] Falha na tentativa %d de conectar: %v\n", consumerName, i+1, err)
		}

		if i < maxRetries-1 {
			time.Sleep(5 * time.Second)
		}
	}

	if db == nil || err != nil {
		fmt.Printf("[%s] Falha ao conectar com PostgreSQL após %d tentativas. Último erro: %v\n", consumerName, maxRetries, err)
		return
	}

	// Configura o driver de migração com timeout
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "transactions",
	})
	if err != nil {
		fmt.Printf("[%s] Erro ao criar a config com o postgres: %v\n", consumerName, err)
		return
	}

	// Cria a instância de migração
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		fmt.Printf("[%s] Erro ao criar a instância de migração: %v\n", consumerName, err)
		return
	}

	// Aplica as migrações com tratamento de erro adequado
	fmt.Printf("[%s] Aplicando migrações...\n", consumerName)
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Printf("[%s] Nenhuma migração pendente\n", consumerName)
		} else {
			fmt.Printf("[%s] Erro ao aplicar as migrações: %v\n", consumerName, err)
			return
		}
	}

	fmt.Printf("[%s] Migrações aplicadas com sucesso\n", consumerName)
}
