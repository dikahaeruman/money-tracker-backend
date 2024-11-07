package repositories

import (
	"context"

	"money-tracker-backend/internal/models"
)

// AccountRepository defines the interface for account-related database operations
type AccountRepository interface {
	// CreateAccount creates a new account in the database
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)

	// GetAccountByID retrieves an account by its ID
	GetAccountByID(ctx context.Context, accountID string) (*models.Account, error)

	// GetAccountsByUserID retrieves all accounts for a given user ID
	GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error)

	// UpdateAccount updates an existing account in the database
	UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error)

	// DeleteAccount removes an account from the database
	DeleteAccount(ctx context.Context, accountID string) error
}
