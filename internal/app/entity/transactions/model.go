package transactions

type Transaction struct {
	ID              uint    `json:"id"`
	TransactionType string  `json:"type"`
	From            uint    `json:"from"`
	To              uint    `json:"to"`
	Amount          float64 `json:"amount"`
	LastTimeUpdated string  `json:"last_time_updated"`
	Comment         string  `json:"comment"`
}
