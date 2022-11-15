package handlersImpl

import (
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *handler) getTransactionsByAccountId(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(res)
	accountId, err := strconv.ParseUint(params["accountId"], 10, 64)
	if err != nil {
		boom.BadData(w, "invalid id type")
		return
	}

	transactions, err := h.transactionRepository.FindByAccountIdUsingStatements(uint(accountId), res.URL.Query())
	if err != nil {
		boom.BadRequest(w, fmt.Sprintf("error while querying transactions for account %d", accountId))
		return
	}

	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		return
	}
}
