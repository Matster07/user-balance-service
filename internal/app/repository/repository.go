package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	Account
	Order
	Service
	Transaction
}

func GetRepository(client *pgxpool.Pool) Repository {
	return Repository{
		Account{client},
		Order{client},
		Service{client},
		Transaction{client},
	}
}
