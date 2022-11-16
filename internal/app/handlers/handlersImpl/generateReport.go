package handlersImpl

import (
	"encoding/csv"
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *handler) generateReport(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(res)
	filename := params["filename"] + ".csv"

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/csv")

	//file, err := os.Create(filename)
	//if err != nil {
	//	boom.BadRequest(w, "error while generating report")
	//	return
	//}

	writer := csv.NewWriter(w)
	err := writer.Write([]string{"service_name", "profit"})
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	rows, err := h.orderRepository.GetDataForReport()
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	err = writer.WriteAll(rows)
	if err != nil {
		return
	}

	writer.Flush()

	//err = json.NewEncoder(w).Encode(rows)
	//if err != nil {
	//	return
	//}
}
