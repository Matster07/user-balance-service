package transactions

import "time"

type TransactionDto struct {
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"type"`
	Comment         string    `json:"comment"`
	CreationDate    time.Time `json:"date"`
}
