package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/app/controller/broker"
	"github.com/matster07/user-balance-service/internal/app/controller/server"
	"github.com/matster07/user-balance-service/internal/app/repository"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"net/http"
	"strconv"
)

func main() {
	cfg := configs.GetConfig()

	logger := logging.GetLogger()
	defer logging.PanicHandler()

	dbClient, err := postgresql.NewClient(context.TODO(), 1, cfg)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	repo := repository.GetRepository(dbClient)

	handler := server.Handler{Repository: repo, Client: dbClient}
	router := mux.NewRouter()
	handler.Register(router)

	reader := broker.GetConsumer()
	reader.Read(handler)

	defer func(pool *pgxpool.Pool, consumer *broker.Consumer) {
		dbClient.Close()

		err := consumer.Reader.Close()
		if err != nil {
			logging.GetLogger().Fatal("failed to close consumer:", err)
		}
	}(dbClient, reader)

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
