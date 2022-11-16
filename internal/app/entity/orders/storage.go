package orders

import "github.com/jackc/pgx/v4"

type Repository interface {
	Create(tx pgx.Tx, order Order) error
	UpdateStatus(tx pgx.Tx, id uint, status string) (order Order, error error)
	GetDataForReport() (result [][]string, err error)
}
