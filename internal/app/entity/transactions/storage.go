package transactions

import (
	"github.com/jackc/pgx/v4"
	"net/url"
)

type Repository interface {
	Create(tx pgx.Tx, model Transaction) error
	FindByAccountIdUsingStatements(accountId uint, query url.Values) (transactions []TransactionDto, err error)
}
