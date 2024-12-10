package models

import (
	"errors"
	"time"
)

type Account struct {
	ID               string    `json:"id"`
	UserID           int       `json:"user_id"`
	AccountName      string    `json:"account_name"`
	Balance          float64   `json:"balance"`
	CurrencyID       int       `json:"currency_id"`
	CurrencyCode     string    `json:"currency_code,omitempty"`
	ConvertedBalance float64   `json:"converted_balance,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// Validate performs validation on the Account struct.
func (a *Account) Validate() error {
	if a.UserID <= 0 {
		return errors.New("user_id is required and must be greater than zero")
	}
	if a.AccountName == "" {
		return errors.New("account_name is required")
	}
	if a.Balance < 0 {
		return errors.New("balance cannot be negative")
	}
	if a.CurrencyID <= 0 {
		return errors.New("currency_id is required")
	}
	return nil
}
