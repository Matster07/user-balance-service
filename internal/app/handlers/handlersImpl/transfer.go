package handlersImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
	"github.com/matster07/user-balance-service/internal/app/entity/transactions"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) transfer(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var transferDto accounts.TransferDTO
	err := json.NewDecoder(res.Body).Decode(&transferDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	from, err := h.accountRepository.FindById(transferDto.From)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", transferDto.From))
		return
	}
	if from.Balance < transferDto.Amount {
		boom.BadRequest(w, "insufficient funds")
		return
	}

	to, err := h.accountRepository.FindById(transferDto.To)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", transferDto.To))
		return
	}

	tx, err := h.client.Begin(context.TODO())

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			h.logger.Trace("Rollback transaction")
		}
	}(tx, context.TODO())

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      transferDto.From,
		Balance: from.Balance - transferDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update account"))
		return
	}

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      transferDto.To,
		Balance: to.Balance + transferDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update account"))
		return
	}

	err = h.transactionRepository.Create(tx, transactions.Transaction{
		Amount:          transferDto.Amount,
		From:            transferDto.From,
		To:              transferDto.To,
		TransactionType: "TRANSFER",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	if err = tx.Commit(context.TODO()); err != nil {
		boom.BadRequest(w, errors.New("failed to commit transaction"))
		return
	}

	returnBalance(w, from.Balance-transferDto.Amount)
}
