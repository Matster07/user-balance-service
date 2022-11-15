package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) FindById(accountId uint) (accounts.Account, error) {
	sql := `
		SELECT id, balance, account_type FROM accounts WHERE id= $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	var acc accounts.Account
	err := r.client.QueryRow(context.TODO(), sql, accountId).Scan(&acc.ID, &acc.Balance, &acc.AccountType)
	if err != nil {
		return accounts.Account{}, err
	}

	return acc, nil
}

func (r *repository) FindByType(accType string) (accounts.Account, error) {
	sql := `
		SELECT id, balance, account_type FROM accounts WHERE account_type = $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	var acc accounts.Account
	err := r.client.QueryRow(context.TODO(), sql, accType).Scan(&acc.ID, &acc.Balance, &acc.AccountType)
	if err != nil {
		return accounts.Account{}, err
	}

	return acc, nil
}

func (r *repository) Create(tx pgx.Tx, accountId uint, balance float64, accountType string) error {
	sql := `
		INSERT INTO accounts 
		    (id, balance, account_type) 
		VALUES 
		       ($1, $2, $3)
		RETURNING ID
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	var accId uint
	if err := tx.QueryRow(context.TODO(), sql, accountId, balance, accountType).Scan(&accId); err != nil {
		return HandlerError(tx, err)
	}

	return nil
}

func (r *repository) Update(tx pgx.Tx, account accounts.Account) error {
	sql := `
		UPDATE accounts 
		SET
		    balance = $2
		WHERE 
			id = $1
		RETURNING ID
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	var accId uint
	if err := tx.QueryRow(context.TODO(), sql, account.ID, account.Balance).Scan(&accId); err != nil {
		return HandlerError(tx, err)
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) accounts.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func HandlerError(tx pgx.Tx, err error) error {
	err1 := tx.Rollback(context.TODO())
	if err1 != nil {
		return err1
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		logging.GetLogger().Error(newErr)
		return newErr
	}
	return err
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
