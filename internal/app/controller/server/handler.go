package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/app/repository"
	"github.com/matster07/user-balance-service/internal/pkg/client"
	"net/http"
)

type Handler struct {
	repository.Repository
	client.Client
}

func (h *Handler) Register(router *mux.Router) {
	prefix := "/api/" + configs.GetConfig().ApiVersion

	// Получение баланса
	router.HandleFunc(prefix+"/accounts/{accountId}/balance", h.getBalanceByAccountId).Methods("GET")

	// Получение транзакций пользователя с пагинацией
	router.HandleFunc(prefix+"/accounts/{accountId}/transactions", h.getTransactionsByAccountId).Methods("GET")

	// Пополнение + инициализация счета
	router.HandleFunc(prefix+"/account/deposit", h.deposit).Methods("POST")

	// Вывод со счета
	router.HandleFunc(prefix+"/account/withdraw", h.withdrawal).Methods("POST")

	// Перевод с одного счета на другой
	router.HandleFunc(prefix+"/account/transfer", h.transfer).Methods("POST")

	// Резервация средств
	router.HandleFunc(prefix+"/account/reserve", h.reserve).Methods("POST")

	// Создание отчета по услугам и прибыли от них
	router.HandleFunc(prefix+"/report/service/profit", h.generateReport).Methods("GET")
}

func returnBalance(w http.ResponseWriter, balance float64) {
	err := json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
	if err != nil {
		return
	}
}

func returnSuccess(w http.ResponseWriter, message string) {
	err := json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		return
	}
}
