package dto

type Account struct {
    AccountName string  `json:"account_name" binding:"required"`
    Balance     float64 `json:"balance" binding:"number,min=0"`
    CurrencyID  int     `json:"currency_id" binding:"required"`
}
