package handlersImpl

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
	"github.com/matster07/user-balance-service/internal/app/entity/categories"
	"github.com/matster07/user-balance-service/internal/app/entity/orders"
	"github.com/matster07/user-balance-service/internal/app/entity/transactions"
	"github.com/matster07/user-balance-service/internal/app/handlers"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"net/http"
)

type handler struct {
	logger                *logging.Logger
	accountRepository     accounts.Repository
	transactionRepository transactions.Repository
	categoryRepository    categories.Repository
	orderRepository       orders.Repository
	client                postgresql.Client
	config                *configs.Config
}

func NewHandler(
	logger *logging.Logger,
	accountRepository accounts.Repository,
	transactionRepository transactions.Repository,
	categoryRepository categories.Repository,
	orderRepository orders.Repository,
	config *configs.Config,
	client postgresql.Client) handlers.Handler {
	return &handler{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		categoryRepository:    categoryRepository,
		orderRepository:       orderRepository,
		logger:                logger,
		config:                config,
		client:                client,
	}
}

func (h *handler) Register(router *mux.Router) {
	prefix := "/api/" + h.config.ApiVersion

	router.HandleFunc(prefix+"/accounts/{accountId}/balance", h.getBalanceByAccountId).Methods("GET")
	router.HandleFunc(prefix+"/accounts/{accountId}/transactions", h.getTransactionsByAccountId).Methods("GET")
	router.HandleFunc(prefix+"/account/deposit", h.deposit).Methods("POST")
	router.HandleFunc(prefix+"/account/withdraw", h.withdrawal).Methods("POST")
	router.HandleFunc(prefix+"/account/transfer", h.transfer).Methods("POST")
	router.HandleFunc(prefix+"/account/reserve", h.reserve).Methods("POST")
	router.HandleFunc(prefix+"/report/service/profit", h.generateReport).Methods("GET")
}

func returnBalance(w http.ResponseWriter, balance float64) {
	err := json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
	if err != nil {
		return
	}
}
