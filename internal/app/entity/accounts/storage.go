package accounts

import "github.com/jackc/pgx/v4"

type Repository interface {
	FindById(accountId uint) (u Account, err error)
	FindByType(accType string) (u Account, err error)
	Create(tx pgx.Tx, accountId uint, balance float64, accountType string) error
	Update(tx pgx.Tx, account Account) error
}
