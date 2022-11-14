package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/accounts/db"
	"github.com/matster07/user-balance-service/internal/app/orders"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(tx pgx.Tx, order orders.Order) error {
	sql := `
		INSERT INTO orders 
		    (id, category_id, price, user_account_id, status) 
		VALUES 
		       ($1, $2, $3, $4, $5)
		RETURNING ID
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	var orderId uint
	if err := tx.QueryRow(context.TODO(), sql, order.ID, order.CategoryId, order.Price, order.UserAccountId, order.Status).Scan(&orderId); err != nil {
		return db.HandlerError(tx, err)
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) orders.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
