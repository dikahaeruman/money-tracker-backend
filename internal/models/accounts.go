package models

type Account struct {
	ID          string  `json:"id"`
	UserID      int     `json:"user_id"`
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	CreatedAt   string  `json:"created_at"`
}
