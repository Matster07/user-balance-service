package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/pkg/errors"
	"net/http"
)

// withdrawal Вывод средств
func (h *Handler) withdrawal(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var accountDto dto.AccDTO
	err := json.NewDecoder(res.Body).Decode(&accountDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	acc, err := h.Account.FindById(accountDto.AccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", accountDto.AccountId))
		return
	}
	if acc.Balance < accountDto.Amount {
		boom.BadRequest(w, "insufficient funds")
		return
	}

	tx, err := h.Begin(context.TODO())

	defer postgresql.RollbackTx(tx)

	err = h.Account.Update(tx, entity.Account{
		ID:      accountDto.AccountId,
		Balance: acc.Balance - accountDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	err = h.Transaction.Create(tx, entity.Transaction{
		Amount:          accountDto.Amount,
		From:            accountDto.AccountId,
		TransactionType: "WITHDRAWAL",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	postgresql.CommitTx(tx)

	returnBalance(w, acc.Balance-accountDto.Amount)
}
