package services

import (
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
)

type AccountService struct {
	accountRepo repositories.AccountRepository
}

func NewAccountService(accountRepo repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (accountService *AccountService) CreateAccount(account *models.Account) (*models.Account, error) {
	return accountService.accountRepo.CreateAccount(account)
}

func (accountService *AccountService) GetAccountByID(accountId string) (*models.Account, error) {
	return accountService.accountRepo.GetAccountByID(accountId)
}

func (accountService *AccountService) GetAccountsByUserID(userId int) ([]*models.Account, error) {
	return accountService.accountRepo.GetAccountsByUserID(userId)
}

func (accountService *AccountService) UpdateAccount(account *models.Account) (*models.Account, error) {
	return accountService.accountRepo.UpdateAccount(account)
}

func (accountService *AccountService) DeleteAccount(accountID string) (string, error) {
	return accountService.accountRepo.DeleteAccount(accountID)
}
