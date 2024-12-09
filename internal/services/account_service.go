package services

import (
	"context"
	"errors"
	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
)

// accountService implements the AccountService interface
type AccountService struct {
	repo interfaces.AccountRepositoryInterface
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(repo interfaces.AccountRepositoryInterface) interfaces.AccountServiceInterface {
	return &AccountService{repo: repo}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	if account == nil {
		return nil, errors.New("account cannot be nil")
	}
	return s.repo.CreateAccount(ctx, account)
}

// GetAccountByID retrieves an account by its ID
func (s *AccountService) GetAccountByID(ctx context.Context, accountID string) (*models.Account, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	return s.repo.GetAccountByID(ctx, accountID)
}

// GetAccounts retrieves all accounts for a given user ID
func (s *AccountService) GetAccounts(ctx context.Context, userID int) ([]*models.Account, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetAccountsByUserID(ctx, userID)
}

// UpdateAccount updates an existing account
func (s *AccountService) UpdateAccount(ctx context.Context, accountID string, accountDTO dto.Account) (*models.Account, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}

	existingAccount, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	existingAccount.AccountName = accountDTO.AccountName
	existingAccount.Balance = accountDTO.Balance
	existingAccount.CurrencyID = accountDTO.CurrencyID

	return s.repo.UpdateAccount(ctx, existingAccount)
}

// DeleteAccount deletes an account by its ID
func (s *AccountService) DeleteAccount(ctx context.Context, accountID string) error {
	if accountID == "" {
		return errors.New("account ID cannot be empty")
	}
	return s.repo.DeleteAccount(ctx, accountID)
}
