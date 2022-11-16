package entity

type Service struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	AccountId   uint   `json:"account_id"`
}
