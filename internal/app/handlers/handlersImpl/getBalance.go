package handlersImpl

import (
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *handler) getBalanceByAccountId(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(res)
	accountId, err := strconv.ParseUint(params["accountId"], 10, 64)
	if err != nil {
		boom.BadData(w, "invalid id type")
		return
	}

	account, err := h.accountRepository.FindById(uint(accountId))
	if err != nil {
		boom.NotFound(w, "accounts id "+params["accountId"]+" wasn't found")
		return
	}

	returnBalance(w, account.Balance)
}
