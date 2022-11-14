package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matster07/user-balance-service/internal/app/configs"
	accountRepo "github.com/matster07/user-balance-service/internal/app/entity/accounts/db"
	categoryRepo "github.com/matster07/user-balance-service/internal/app/entity/categories/db"
	orderRepo "github.com/matster07/user-balance-service/internal/app/entity/orders/db"
	transactionRepo "github.com/matster07/user-balance-service/internal/app/entity/transactions/db"
	"github.com/matster07/user-balance-service/internal/app/handlers/handlersImpl"
	"github.com/matster07/user-balance-service/internal/app/kafka/consumer"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"net/http"
	"strconv"
)

func main() {
	cfg := configs.GetConfig()

	logger := logging.GetLogger()
	dbClient, err := postgresql.NewClient(context.TODO(), 1, cfg)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	reader := consumer.GetConsumer()
	reader.Read()

	defer func(pool *pgxpool.Pool, consumer *consumer.Consumer) {
		dbClient.Close()

		err := consumer.Reader.Close()
		if err != nil {
			logging.GetLogger().Fatal(err)
		}
	}(dbClient, reader)

	router := mux.NewRouter()

	accountRepository := accountRepo.NewRepository(dbClient, logger)
	transactionRepository := transactionRepo.NewRepository(dbClient, logger)
	categoryRepository := categoryRepo.NewRepository(dbClient, logger)
	orderRepository := orderRepo.NewRepository(dbClient, logger)

	handler := handlersImpl.NewHandler(logger, accountRepository, transactionRepository, categoryRepository, orderRepository, cfg, dbClient)
	handler.Register(router)

	startServer(router, cfg)
}

func startServer(router *mux.Router, cfg *configs.Config) {
	logger := logging.GetLogger()

	logger.Infof("server is listening port %d", cfg.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
	if err != nil {
		panic(err)
	}
}
