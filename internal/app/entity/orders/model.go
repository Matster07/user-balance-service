package orders

import "time"

type Order struct {
	ID            uint      `json:"id"`
	CategoryId    uint      `json:"category_id"`
	Price         float64   `json:"price"`
	UserAccountId uint      `json:"user_account_id"`
	CreationDate  time.Time `json:"creation_date"`
	Status        string    `json:"status"`
}
