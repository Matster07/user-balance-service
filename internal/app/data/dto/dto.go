package dto

type AccDTO struct {
	AccountId uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type TransferDTO struct {
	From   uint    `json:"from"`
	To     uint    `json:"to"`
	Amount float64 `json:"amount"`
}

type ReserveDTO struct {
	AccountId uint    `json:"user_account_id"`
	ServiceId uint    `json:"service_id"`
	OrderId   uint    `json:"order_id"`
	Price     float64 `json:"price"`
}

type DeliverStatusDto struct {
	Status        string  `json:"status"`
	OrderId       uint    `json:"order_id"`
	UserAccountId uint    `json:"user_account_id"`
	Amount        float64 `json:"amount"`
	ServiceId     uint    `json:"service_id"`
}
