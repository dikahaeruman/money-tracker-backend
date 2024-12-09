package dto

type Account struct {
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
	CurrencyID  int     `json:"currency_id"`
}
