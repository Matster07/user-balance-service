package server

import (
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//	@Summary      Get account transactions history
//	@Description  Получение истории транзакций счета с пагинацией, фильтрацией
//	@Tags         account
//	@Accept       json
//	@Produce      json
//  @Param        accountId      path  uint   true  "Идентификатор счета"
//  @Param        amount_sort    query string false "Сортировка по сумме транзакции(asc/desc)"
//  @Param        date_sort      query string false "Сортировка по дате (asc/desc)"
//  @Param        page           query uint false "Пагинация по странице. Страница вмешает 9 значений"
//  @Success      200            {object} []entity.TransactionPagination
//	@Router       /accounts/{accountId}/transactions [get]]
func (h *Handler) getTransactionsByAccountId(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(res)
	accountId, err := strconv.ParseUint(params["accountId"], 10, 64)
	if err != nil {
		boom.BadData(w, "invalid id type")
		return
	}

	transactions, err := h.Transaction.FindByAccountIdUsingStatements(uint(accountId), res.URL.Query())
	if err != nil {
		boom.BadRequest(w, fmt.Sprintf("error while querying transactions for account %d", accountId))
		return
	}

	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		return
	}
}
