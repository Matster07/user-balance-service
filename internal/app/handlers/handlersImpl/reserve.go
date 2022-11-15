package handlersImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
	"github.com/matster07/user-balance-service/internal/app/entity/orders"
	"github.com/matster07/user-balance-service/internal/app/entity/transactions"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) reserve(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reserveDto accounts.ReserveDTO
	err := json.NewDecoder(res.Body).Decode(&reserveDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	account, err := h.accountRepository.FindById(reserveDto.AccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", reserveDto.AccountId))
		return
	}
	if account.Balance < reserveDto.Price {
		boom.BadRequest(w, "insufficient funds")
		return
	}

	category, err := h.categoryRepository.FindById(reserveDto.CategoryId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("category %d wasn't found", reserveDto.CategoryId))
		return
	}
	reservingAccountId := category.AccountId
	reservingAcc, err := h.accountRepository.FindById(reservingAccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("reserving account %d wasn't found", reserveDto.AccountId))
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
		ID:      account.ID,
		Balance: account.Balance - reserveDto.Price,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update balance"))
		return
	}

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      reservingAcc.ID,
		Balance: reservingAcc.Balance + reserveDto.Price,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update balance"))
		return
	}

	order := orders.Order{
		ID:            reserveDto.OrderId,
		UserAccountId: account.ID,
		Price:         reserveDto.Price,
		CategoryId:    category.ID,
		Status:        "IN_PROGRESS",
	}
	err = h.orderRepository.Create(tx, order)
	if err != nil {
		boom.BadRequest(w, errors.New("failed to process order"))
		return
	}

	err = h.transactionRepository.Create(tx, transactions.Transaction{
		Amount:          reserveDto.Price,
		From:            account.ID,
		To:              reservingAcc.ID,
		TransactionType: "RESERVE",
		Comment:         fmt.Sprintf("order_id: %d", order.ID),
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	if err = tx.Commit(context.TODO()); err != nil {
		boom.BadRequest(w, errors.New("failed to commit transaction"))
		return
	}
}
