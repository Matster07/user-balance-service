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

// reserve Резервирование средств
func (h *Handler) reserve(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reserveDto dto.ReserveDTO
	err := json.NewDecoder(res.Body).Decode(&reserveDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	account, err := h.Account.FindById(reserveDto.AccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", reserveDto.AccountId))
		return
	}
	if account.Balance < reserveDto.Price {
		boom.BadRequest(w, "insufficient funds")
		return
	}

	service, err := h.Service.FindById(reserveDto.ServiceId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("service %d wasn't found", reserveDto.ServiceId))
		return
	}
	reservingAccountId := service.AccountId
	reservingAcc, err := h.Account.FindById(reservingAccountId)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("reserving account %d wasn't found", reserveDto.AccountId))
		return
	}

	tx, err := h.Begin(context.TODO())

	defer postgresql.RollbackTx(tx)

	err = h.Account.Update(tx, entity.Account{
		ID:      account.ID,
		Balance: account.Balance - reserveDto.Price,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update balance"))
		return
	}

	err = h.Account.Update(tx, entity.Account{
		ID:      reservingAcc.ID,
		Balance: reservingAcc.Balance + reserveDto.Price,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update balance"))
		return
	}

	order := entity.Order{
		ID:            reserveDto.OrderId,
		UserAccountId: account.ID,
		Price:         reserveDto.Price,
		ServiceId:     service.ID,
		Status:        "IN_PROGRESS",
	}
	err = h.Order.Create(tx, order)
	if err != nil {
		boom.BadRequest(w, errors.New("failed to process order"))
		return
	}

	err = h.Transaction.Create(tx, entity.Transaction{
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

	postgresql.CommitTx(tx)

	returnSuccess(w, "SUCCESS")
}
