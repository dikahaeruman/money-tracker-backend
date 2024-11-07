package services

import (
	"context"
	"errors"
	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
)

// AccountService defines the interface for account-related operations
type AccountService interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*models.Account, error)
	GetAccounts(ctx context.Context, userID int) ([]*models.Account, error)
	UpdateAccount(ctx context.Context, accountID string, accountDTO dto.Account) (*models.Account, error)
	DeleteAccount(ctx context.Context, accountID string) error
}

// accountService implements the AccountService interface
type accountService struct {
	repo repositories.AccountRepository
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// CreateAccount creates a new account
func (s *accountService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	if account == nil {
		return nil, errors.New("account cannot be nil")
	}
	return s.repo.CreateAccount(ctx, account)
}

// GetAccountByID retrieves an account by its ID
func (s *accountService) GetAccountByID(ctx context.Context, accountID string) (*models.Account, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	return s.repo.GetAccountByID(ctx, accountID)
}

// GetAccounts retrieves all accounts for a given user ID
func (s *accountService) GetAccounts(ctx context.Context, userID int) ([]*models.Account, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetAccountsByUserID(ctx, userID)
}

// UpdateAccount updates an existing account
func (s *accountService) UpdateAccount(ctx context.Context, accountID string, accountDTO dto.Account) (*models.Account, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}

	existingAccount, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	existingAccount.AccountName = accountDTO.AccountName
	existingAccount.Balance = accountDTO.Balance
	existingAccount.Currency = accountDTO.Currency

	return s.repo.UpdateAccount(ctx, existingAccount)
}

// DeleteAccount deletes an account by its ID
func (s *accountService) DeleteAccount(ctx context.Context, accountID string) error {
	if accountID == "" {
		return errors.New("account ID cannot be empty")
	}
	return s.repo.DeleteAccount(ctx, accountID)
}
