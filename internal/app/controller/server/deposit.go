package server

import (
	"context"
	"encoding/json"
	"github.com/darahayes/go-boom"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/pkg/errors"
	"net/http"
)

// deposit Пополнение счета
func (h *Handler) deposit(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var accountDto dto.AccDTO
	err := json.NewDecoder(res.Body).Decode(&accountDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	tx, err := h.Begin(context.TODO())

	defer postgresql.RollbackTx(tx)

	acc, err := h.Account.FindById(accountDto.AccountId)
	balance := accountDto.Amount + acc.Balance
	if err != nil {
		if err = h.Account.Create(tx, accountDto.AccountId, accountDto.Amount, "CUSTOMER"); err != nil {
			boom.BadData(w, errors.New("failed to process transaction"))
			return
		}
	}

	err = h.Account.Update(tx, entity.Account{
		ID:      accountDto.AccountId,
		Balance: balance,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update account"))
		return
	}

	err = h.Transaction.Create(tx, entity.Transaction{
		Amount:          accountDto.Amount,
		To:              accountDto.AccountId,
		TransactionType: "DEPOSIT",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	postgresql.CommitTx(tx)

	returnBalance(w, balance)
}
