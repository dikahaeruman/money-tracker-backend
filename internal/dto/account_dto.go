package dto

type Account struct {
	AccountName string  `json:"account_name" binding:"required"`
	Balance     float64 `json:"balance" binding:"required"`
	CurrencyID  int     `json:"currency_id" binding:"required"`
}
