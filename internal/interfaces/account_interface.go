package interfaces

import (
	"context"
	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// AccountRepository defines the interface for account-related database operations
type AccountRepositoryInterface interface {
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

// AccountService defines the interface for account-related operations
type AccountServiceInterface interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*models.Account, error)
	GetAccounts(ctx context.Context, userID int) ([]*models.Account, error)
	UpdateAccount(ctx context.Context, accountID string, accountDTO dto.Account) (*models.Account, error)
	DeleteAccount(ctx context.Context, accountID string) error
}

type AccountControllerInterface interface {
	CreateAccount(c *gin.Context)
	GetAccountByID(c *gin.Context)
	GetAccounts(c *gin.Context)
	UpdateAccount(c *gin.Context)
	DeleteAccount(c *gin.Context)
}
