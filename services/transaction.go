package services

import (
	"context"
	"database/sql"
	"fmt"
	"main/entities"
	"main/pkg/database"
	"os"

	"github.com/uptrace/bun"
)

const (
	OperacaoCredito = "c"
	OperacaoDebito  = "d"
)

func FindClient(ctx context.Context, db *bun.DB, id string) (*entities.Client, error) {
	functionName := "FindClient"
	cliente := new(entities.Client)

	err := db.NewSelect().Model(cliente).Where("id_client = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return cliente, err
	}
	return cliente, nil
}

func FindClientTx(ctx context.Context, tx bun.Tx, id string) (*entities.Client, error) {
	functionName := "FindClient"
	cliente := new(entities.Client)

	err := tx.NewSelect().Model(cliente).Where("id_client = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Error to recover client %v:\n", functionName, err)
		tx.Rollback()
		return cliente, err
	}
	return cliente, nil
}

func Process(transaction entities.Transaction) (novoBalance int64, limit int64, inconsistency bool, err error) {
	functionName := fmt.Sprintf("OperacaoDe%s", transaction.Type)

	ctx := context.Background()
	db := database.GetDB()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		fmt.Printf("[%s] Error to init database transaction: %v\n", functionName, err)
		return
	}

	// Garantir que um rollback será feito em caso de erro
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	client, err := FindClientTx(ctx, tx, transaction.IDClient)
	if err != nil {
		fmt.Printf("[%s] Error to find clients %v:\n", functionName, err)
		return
	}

	// Cálculo do novo balance com base no type de operação
	switch transaction.Type {
	case OperacaoCredito:
		novoBalance = client.Balance + transaction.Amount
	case OperacaoDebito:
		novoBalance = client.Balance - transaction.Amount
		if novoBalance < -client.Limit {
			inconsistency = true
			err = fmt.Errorf("[%s] No limit for client", functionName)
			return 0, 0, true, err
		}
	default:
		err = fmt.Errorf("[%s] Invalid transaction type", functionName)
		return 0, 0, false, err
	}

	// Atualizar balance do cliente
	// _, err = tx.ExecContext(ctx, "UPDATE clients SET balance = ? WHERE id_client = ?", novoBalance, transactionIDClient)
	// if err != nil {
	// 	fmt.Printf("[%s] Erro ao atualizar o balance do cliente: %v\n", functionName, err)
	// 	return
	// }

	_, err = tx.NewUpdate().Model((*entities.Client)(nil)).
		Set("balance = ?", novoBalance).
		Where("id_client = ?", transaction.IDClient).
		Exec(ctx)

	if err != nil {
		fmt.Printf("[%s] Erro ao atualizar o balance do cliente: %v\n", functionName, err)
		return 0, 0, false, err
	}

	_, err = tx.NewInsert().
		Model(&transaction).
		Exec(ctx)

	// Inserir transação
	// _, err = tx.ExecContext(ctx, "INSERT INTO transactions (id_client, amount, type, description) VALUES (?, ?, ?, ?)", transactionIDClient, transactionAmount, type, transactionDescription)
	if err != nil {
		fmt.Printf("[%s] Erro ao inserir a transação: %v\n", functionName, err)
		return 0, 0, false, err
	}

	if os.Getenv("ENV") == "shadow" {
		err = tx.Rollback()
		if err != nil {
			fmt.Printf("[%s] Erro ao fazer rollback da transação: %v\n", functionName, err)
			return 0, 0, false, err
		}
		return novoBalance, client.Limit, inconsistency, nil
	}

	// Commit da Transação
	err = tx.Commit()
	if err != nil {
		fmt.Printf("[%s] Erro ao fazer commit da transação: %v\n", functionName, err)
		return 0, 0, false, err
	}

	return novoBalance, client.Limit, inconsistency, nil
}
