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

func (h *handler) withdrawal(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var accountDto accounts.AccDTO
	err := json.NewDecoder(res.Body).Decode(&accountDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	acc, err := h.accountRepository.FindById(accountDto.AccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", accountDto.AccountId))
		return
	}
	if acc.Balance < accountDto.Amount {
		boom.BadRequest(w, "insufficient funds")
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
		ID:      accountDto.AccountId,
		Balance: acc.Balance - accountDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	err = h.transactionRepository.Create(tx, transactions.Transaction{
		Amount:          accountDto.Amount,
		From:            accountDto.AccountId,
		TransactionType: "WITHDRAWAL",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	if err = tx.Commit(context.TODO()); err != nil {
		boom.BadRequest(w, errors.New("failed to commit transaction"))
		return
	}

	returnBalance(w, acc.Balance-accountDto.Amount)
}
