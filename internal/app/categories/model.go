package categories

type Category struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"service_name"`
	AccountId    uint   `json:"account_id"`
}
