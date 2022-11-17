package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"net/http"
	"os"
	filepath2 "path/filepath"
	"strconv"
)

// generateReport Генерация отчет по всем пользователям с указанием выручки по каждой услуге
func (h *Handler) generateReport(w http.ResponseWriter, res *http.Request) {
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

	rows, err := h.Order.GetDataForReport(uint(year), uint(month))
	if err != nil {
		boom.BadRequest(w, "error while generating report")
		return
	}

	h.populateZeroProfitServices(&rows)

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

func (h *Handler) populateZeroProfitServices(result *[][]string) {
	services, err := h.Service.FindAll()
	if err != nil {
		logging.GetLogger().Errorf("error while fetching services")
	}

	for _, value := range services {
		needToAdd := true

		for _, row := range *result {
			if row[0] == value.ServiceName {
				needToAdd = false
			}
		}

		if needToAdd {
			*result = append(*result, []string{value.ServiceName, "0"})
		}
	}
}
