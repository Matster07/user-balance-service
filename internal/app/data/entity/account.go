package entity

type Account struct {
	ID          uint    `json:"id"`
	Balance     float64 `json:"balance"`
	AccountType string  `json:"account_type"`
}
