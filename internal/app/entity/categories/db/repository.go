package db

import (
	"context"
	"fmt"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts/db"
	"github.com/matster07/user-balance-service/internal/app/entity/categories"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) FindById(id uint) (category categories.Category, err error) {
	sql := `
		SELECT id, account_id, category_name FROM categories WHERE id= $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	if err = r.client.QueryRow(context.TODO(), sql, id).Scan(&category.ID, &category.AccountId, &category.CategoryName); err != nil {
		return categories.Category{}, err
	}

	return category, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) categories.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
