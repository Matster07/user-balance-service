package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts/db"
	"github.com/matster07/user-balance-service/internal/app/entity/orders"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) UpdateStatus(tx pgx.Tx, id uint, status string) (order orders.Order, err error) {
	sql := `
		UPDATE orders 
		SET
		    status = $2
		WHERE 
			id = $1
		RETURNING ID, price, category_id, status
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	err = tx.QueryRow(context.TODO(), sql, id, status).Scan(&order.ID, &order.Price, &order.CategoryId, &order.Status)
	if err != nil {
		return order, err
	}

	return order, nil
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

func (r *repository) GetDataForReport() (result [][]string, err error) {
	sql := `
		SELECT category_name, SUM(price) AS profit FROM orders JOIN categories c on c.id = orders.category_id GROUP BY category_name
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	rows, err := r.client.Query(context.TODO(), sql)
	if err != nil {
		logging.GetLogger().Errorf(err.Error())
		return result, err
	}

	result = make([][]string, 0)

	for rows.Next() {
		var categoryName string
		var profit string

		err = rows.Scan(&categoryName, &profit)
		if err != nil {
			return nil, err
		}

		result = append(result, []string{categoryName, profit})
	}

	return result, err
}

//
//type reportRow []string
//
//func (r *reportRow) append() {
//
//}

func NewRepository(client postgresql.Client, logger *logging.Logger) orders.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
