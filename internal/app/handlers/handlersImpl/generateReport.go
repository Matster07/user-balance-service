package handlersImpl

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"net/http"
	"os"
	filepath2 "path/filepath"
	"strconv"
)

func (h *handler) generateReport(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	year, err := strconv.ParseUint(res.URL.Query().Get("year"), 10, 16)
	month, err := strconv.ParseUint(res.URL.Query().Get("month"), 10, 16)

	filename := filepath2.Join("reports", fmt.Sprintf("%d-%d-profit-report.csv", year, month))

	file, err := os.Create(filename)
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	writer := csv.NewWriter(file)
	err = writer.Write([]string{"service_name", "profit"})
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	rows, err := h.orderRepository.GetDataForReport(uint(year), uint(month))
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	err = writer.WriteAll(rows)
	if err != nil {
		return
	}

	writer.Flush()

	err = json.NewEncoder(w).Encode(map[string]string{"stats": "success"})
	if err != nil {
		return
	}
}
