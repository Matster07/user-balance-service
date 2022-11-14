package db

import (
	"context"
	"fmt"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts/db"
	categories2 "github.com/matster07/user-balance-service/internal/app/entity/categories"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) FindById(id uint) (category categories2.Category, err error) {
	sql := `
		SELECT id, account_id, category_name FROM categories WHERE id= $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", db.FormatQuery(sql)))

	if err = r.client.QueryRow(context.TODO(), sql, id).Scan(&category.ID, &category.AccountId, &category.CategoryName); err != nil {
		return categories2.Category{}, err
	}

	return category, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) categories2.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
