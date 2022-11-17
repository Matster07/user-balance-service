package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	_ "github.com/matster07/user-balance-service/docs"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
	"github.com/pkg/errors"
	"net/http"
)

//	@Summary      Transfer
//	@Description  Перевод средств с одного указанного счета на другой
//	@Tags         account
//	@Accept       json
//	@Produce      json
//  @Param        TransferDTO body dto.TransferDTO  true "Идентификатор счета отправителя, идентификатор счета получетеля, сумма перевода"
//	@Router       /account/transfer [post]
func (h *Handler) transfer(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var transferDto dto.TransferDTO
	err := json.NewDecoder(res.Body).Decode(&transferDto)
	if err != nil {
		boom.BadData(w, "invalid body format")
		return
	}

	from, err := h.Account.FindById(transferDto.From)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", transferDto.From))
		return
	}
	if from.Balance < transferDto.Amount {
		boom.BadRequest(w, "insufficient funds")
		return
	}

	to, err := h.Account.FindById(transferDto.To)
	if err != nil {
		boom.NotFound(w, fmt.Sprintf("accounts id %d wasn't found", transferDto.To))
		return
	}

	tx, err := h.Begin(context.TODO())

	defer postgresql.RollbackTx(tx)

	err = h.Account.Update(tx, entity.Account{
		ID:      transferDto.From,
		Balance: from.Balance - transferDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update account"))
		return
	}

	err = h.Account.Update(tx, entity.Account{
		ID:      transferDto.To,
		Balance: to.Balance + transferDto.Amount,
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to update account"))
		return
	}

	err = h.Transaction.Create(tx, entity.Transaction{
		Amount:          transferDto.Amount,
		From:            transferDto.From,
		To:              transferDto.To,
		TransactionType: "TRANSFER",
	})
	if err != nil {
		boom.BadRequest(w, errors.New("failed to save transaction"))
		return
	}

	postgresql.CommitTx(tx)

	dto.ReturnStatus(w, "SUCCESS")
}
