package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"github.com/pkg/errors"
)

type Order struct {
	client.Client
}

func (r *Order) UpdateStatus(tx pgx.Tx, id uint, status string) (err error) {
	sql := `
		UPDATE orders 
		SET
		    status = $2
		WHERE 
			id = $1
		RETURNING ID
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	err = tx.QueryRow(context.TODO(), sql, id, status).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Order) FindById(id uint) (order entity.Order, err error) {
	sql := `
		SELECT id, service_id, price, user_account_id, status, creation_date
		FROM orders WHERE id = $1
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	err = r.QueryRow(context.TODO(), sql, id).Scan(&order.ID, &order.ServiceId, &order.Price, &order.UserAccountId, &order.Status, &order.CreationDate)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *Order) Create(tx pgx.Tx, order entity.Order) error {
	sql := `
		INSERT INTO orders 
		    (id, service_id, price, user_account_id, status) 
		VALUES 
		       ($1, $2, $3, $4, $5)
		RETURNING ID
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	var orderId uint
	if err := tx.QueryRow(context.TODO(), sql, order.ID, order.ServiceId, order.Price, order.UserAccountId, order.Status).Scan(&orderId); err != nil {
		return HandlerError(tx, err)
	}

	return nil
}

func (r *Order) GetDataForReport(year uint, month uint) (result [][]string, err error) {
	if month > 12 {
		return result, errors.New("incorrect month value")
	}

	sql := `
		SELECT service_name, SUM(price) AS total_profit 
        FROM orders
		JOIN services s on s.id = orders.service_id
        WHERE EXTRACT(YEAR FROM creation_date) = $1 AND EXTRACT(MONTH FROM creation_date) = $2 AND status = 'COMPLETED'
        GROUP BY service_name
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	rows, err := r.Query(context.TODO(), sql, year, month)
	if err != nil {
		logging.GetLogger().Errorf(err.Error())
		return result, err
	}

	result = make([][]string, 0)

	for rows.Next() {
		var serviceName string
		var profit string

		err = rows.Scan(&serviceName, &profit)
		if err != nil {
			return nil, err
		}

		result = append(result, []string{serviceName, profit})
	}

	return result, err
}
