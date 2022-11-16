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
