package dto

type Account struct {
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
}
