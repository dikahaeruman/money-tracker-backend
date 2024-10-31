package repositories

import (
	"money-tracker-backend/internal/models"
)

type AccountRepository interface {
	CreateAccount(account *models.Account) (*models.Account, error)
	GetAccountByID(accountID string) (*models.Account, error)
	GetAccountsByUserID(userID int) ([]*models.Account, error)
	UpdateAccount(account *models.Account) (*models.Account, error)
	DeleteAccount(accountID string) (string, error)
}
