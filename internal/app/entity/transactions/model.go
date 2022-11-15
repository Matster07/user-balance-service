package transactions

import (
	"time"
)

type Transaction struct {
	ID              uint      `json:"id"`
	TransactionType string    `json:"type"`
	From            uint      `json:"sender_id"`
	To              uint      `json:"receiver_id"`
	Amount          float64   `json:"amount"`
	CreationDate    time.Time `json:"creation_date"`
	Comment         string    `json:"comment"`
}
