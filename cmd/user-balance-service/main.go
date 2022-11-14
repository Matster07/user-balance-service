package main

import (
	"context"
	"github.com/gorilla/mux"
	accountRepo "github.com/matster07/user-balance-service/internal/app/accounts/db"
	categoryRepo "github.com/matster07/user-balance-service/internal/app/categories/db"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/app/handlers/handlersImpl"
	orderRepo "github.com/matster07/user-balance-service/internal/app/orders/db"
	transactionRepo "github.com/matster07/user-balance-service/internal/app/transactions/db"
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

	defer dbClient.Close()

	router := mux.NewRouter()

	accountRepository := accountRepo.NewRepository(dbClient, logger)
	transactionRepository := transactionRepo.NewRepository(dbClient, logger)
	categoryRepository := categoryRepo.NewRepository(dbClient, logger)
	orderRepository := orderRepo.NewRepository(dbClient, logger)

	handler := handlersImpl.NewHandler(logger, accountRepository, transactionRepository, categoryRepository, orderRepository, cfg, dbClient)
	handler.Register(router)

	start(router, cfg)
}

func start(router *mux.Router, cfg *configs.Config) {
	logger := logging.GetLogger()

	logger.Infof("server is listening port %d", cfg.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
	if err != nil {
		panic(err)
	}
}
