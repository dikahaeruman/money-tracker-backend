package services

import (
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (accountService *AccountService) GetAllAccounts(userId string) ([]*models.Account, error) {
	return accountService.accountRepo.GetAll(userId)
}

func (accountService *AccountService) GetAccount(accountId string) (*models.Account, error) {
	return accountService.accountRepo.Get(accountId)
}

func (accountService *AccountService) CreateAccount(account *models.Account) (*models.Account, error) {
	return accountService.accountRepo.Create(account)
}
