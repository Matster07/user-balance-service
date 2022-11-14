package handlersImpl

import (
	"context"
	"encoding/json"
	"github.com/darahayes/go-boom"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/accounts"
	"github.com/matster07/user-balance-service/internal/app/transactions"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) deposit(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var accountDto accounts.AccDTO
	err := json.NewDecoder(res.Body).Decode(&accountDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	tx, err := h.client.Begin(context.TODO())

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			h.logger.Trace("Rollback transaction")
		}
	}(tx, context.TODO())

	acc, err := h.accountRepository.FindById(accountDto.AccountId)
	balance := accountDto.Amount + acc.Balance
	if err != nil {
		if err = h.accountRepository.Create(tx, accountDto.AccountId, accountDto.Amount, "CUSTOMER"); err != nil {
			boom.BadData(w, errors.New("failed to process transaction"))
			return
		}
	}

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      accountDto.AccountId,
		Balance: balance,
	})
	if err != nil {
		boom.BadData(w, errors.New(err.Error()))
		return
	}

	err = h.transactionRepository.Create(tx, transactions.Transaction{
		Amount:          accountDto.Amount,
		To:              accountDto.AccountId,
		TransactionType: "DEPOSIT",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	if err = tx.Commit(context.TODO()); err != nil {
		boom.BadRequest(w, errors.New("failed to commit transaction"))
		return
	}

	returnBalance(w, balance)
}
