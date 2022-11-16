package entity

import "time"

type Order struct {
	ID            uint      `json:"id"`
	ServiceId     uint      `json:"service_id"`
	Price         float64   `json:"price"`
	UserAccountId uint      `json:"user_account_id"`
	CreationDate  time.Time `json:"creation_date"`
	Status        string    `json:"status"`
}
