package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts/db"
	"github.com/matster07/user-balance-service/internal/app/entity/transactions"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"net/url"
	"strconv"
)

type repository struct {
	logger *logging.Logger
	client postgresql.Client
}

func (r *repository) FindByAccountIdUsingStatements(accountId uint, query url.Values) (trans []transactions.TransactionDto, err error) {
	sql := `
		SELECT type, amount, creation_date, comment FROM transactions WHERE sender_id = $1 OR receiver_id = $1
	`

	if amountSort := query.Get("amount_sort"); amountSort != "" {
		sql = fmt.Sprintf("%s ORDER BY amount %s", sql, amountSort)
	} else if dateSort := query.Get("date_sort"); dateSort != "" {
		sql = fmt.Sprintf("%s ORDER BY creation_date %s", sql, dateSort)
	}

	page, err := strconv.Atoi(query.Get("page"))
	if page != 0 && err == nil {
		perPage := 9
		sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perPage, (page-1)*perPage)
	}
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	rows, err := r.client.Query(context.TODO(), sql, accountId)
	if err != nil {
		logging.GetLogger().Errorf(err.Error())
		return trans, err
	}

	trans = make([]transactions.TransactionDto, 0)

	for rows.Next() {
		transaction := transactions.TransactionDto{}

		err = rows.Scan(
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.CreationDate,
			&transaction.Comment)
		if err != nil {
			return trans, err
		}

		trans = append(trans, transaction)
	}

	if err = rows.Err(); err != nil {
		return trans, err
	}

	return trans, err
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
