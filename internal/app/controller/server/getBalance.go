package server

import (
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	_ "github.com/matster07/user-balance-service/docs"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"net/http"
	"strconv"
)

//	@Summary      Get account balance
//	@Description  Получение баланса счета по его идентификатору
//	@Tags         account
//	@Accept       json
//	@Produce      json
//  @Param        accountId path uint true "Идентификатор счета"
//  @Success      200            {object} dto.BalanceDTO
//	@Router       /accounts/{accountId}/balance [get]
func (h *Handler) getBalanceByAccountId(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(res)
	accountId, err := strconv.ParseUint(params["accountId"], 10, 64)
	if err != nil {
		boom.BadData(w, "invalid id type")
		return
	}

	account, err := h.Account.FindById(uint(accountId))
	if err != nil {
		boom.NotFound(w, "accounts id "+params["accountId"]+" wasn't found")
		return
	}

	dto.ReturnBalance(w, account.Balance)
}
