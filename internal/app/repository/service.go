package repository

import (
	"context"
	"fmt"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type Service struct {
	client.Client
}

func (r *Service) FindById(id uint) (service entity.Service, err error) {
	sql := `
		SELECT id, account_id, service_name FROM services WHERE id= $1;
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	if err = r.QueryRow(context.TODO(), sql, id).Scan(&service.ID, &service.AccountId, &service.ServiceName); err != nil {
		return entity.Service{}, err
	}

	return service, nil
}

func (r *Service) FindAll() (service []entity.Service, err error) {
	sql := `
		SELECT id, account_id, service_name FROM services;
	`
	logging.GetLogger().Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(sql)))

	rows, err := r.Query(context.TODO(), sql)
	if err != nil {
		return []entity.Service{}, err
	}

	service = make([]entity.Service, 0)

	for rows.Next() {
		var row entity.Service

		if err = rows.Scan(&row.ID, &row.AccountId, &row.ServiceName); err != nil {
			return nil, err
		}

		service = append(service, row)
	}

	return service, nil
}
