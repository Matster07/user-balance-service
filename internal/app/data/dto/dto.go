package dto

import (
	"encoding/json"
	"net/http"
)

type AccDTO struct {
	AccountId uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type TransferDTO struct {
	From   uint    `json:"from"`
	To     uint    `json:"to"`
	Amount float64 `json:"amount"`
}

type ReserveDTO struct {
	AccountId uint    `json:"user_account_id"`
	ServiceId uint    `json:"service_id"`
	OrderId   uint    `json:"order_id"`
	Price     float64 `json:"price"`
}

type DeliverStatusDTO struct {
	Status        string  `json:"status"`
	OrderId       uint    `json:"order_id"`
	UserAccountId uint    `json:"user_account_id"`
	Amount        float64 `json:"amount"`
	ServiceId     uint    `json:"service_id"`
}

type BalanceDTO struct {
	Balance float64 `json:"balance"`
}

func ReturnBalance(w http.ResponseWriter, amount float64) {
	err := json.NewEncoder(w).Encode(BalanceDTO{amount})
	if err != nil {
		return
	}
}

func ReturnStatus(w http.ResponseWriter, message string) {
	err := json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		return
	}
}
