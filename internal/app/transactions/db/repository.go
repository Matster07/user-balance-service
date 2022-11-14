package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/accounts/db"
	"github.com/matster07/user-balance-service/internal/app/transactions"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type repository struct {
	logger *logging.Logger
	client postgresql.Client
}

func (r *repository) Create(tx pgx.Tx, model transactions.Transaction) error {
	sql := `
		INSERT INTO transactions 
		    (type, amount, sender_id, receiver_id, comment) 
		VALUES 
		       ($1, $2, $3, $4, $5)
		RETURNING ID
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	var transactionId uint
	if err := prepareArgsByTransactionType(tx, model, sql).Scan(&transactionId); err != nil {
		return db.HandlerError(tx, err)
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) transactions.Repository {
	return &repository{
		logger: logger,
		client: client,
	}
}

func prepareArgsByTransactionType(tx pgx.Tx, t transactions.Transaction, sql string) pgx.Row {
	switch t.TransactionType {
	case "DEPOSIT":
		return tx.QueryRow(context.TODO(), sql, t.TransactionType, t.Amount, nil, t.To, t.Comment)
	case "WITHDRAWAL":
		return tx.QueryRow(context.TODO(), sql, t.TransactionType, t.Amount, t.From, nil, t.Comment)
	default:
		return tx.QueryRow(context.TODO(), sql, t.TransactionType, t.Amount, t.From, t.To, t.Comment)
	}
}
