package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"github.com/matster07/user-balance-service/internal/pkg/utils"
	"time"
)

func NewClient(ctx context.Context, maxAttempts int, sc *configs.Config) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?search_path=%s",
		sc.DatabaseUsername,
		sc.DatabasePassword,
		sc.DatabaseHost,
		sc.DatabasePort,
		sc.DatabaseTable,
		sc.DatabaseSchema)

	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logging.GetLogger().Fatal(err)
	}

	return pool, nil
}

func RollbackTx(tx pgx.Tx) {
	err := tx.Rollback(context.TODO())
	if err != nil {
		logging.GetLogger().Trace("Rollback transaction")
	}
}

func CommitTx(tx pgx.Tx) {
	err := tx.Commit(context.TODO())
	if err != nil {
		logging.GetLogger().Trace("Commit transaction")
	}
}
