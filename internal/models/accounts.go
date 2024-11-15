package models

import (
	"errors"
	"time"
)

type Account struct {
	ID          string    `json:"id"`
	UserID      int       `json:"user_id"`
	AccountName string    `json:"account_name"`
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
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
	if a.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}
