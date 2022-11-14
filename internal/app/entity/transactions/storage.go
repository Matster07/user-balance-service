package transactions

import "github.com/jackc/pgx/v4"

type Repository interface {
	Create(tx pgx.Tx, model Transaction) error
}
