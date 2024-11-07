package services

import (
	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
)

type AccountService struct {
	accountRepository repositories.AccountRepository
}

func NewAccountService(accountRepo repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepository: accountRepo}
}

func (accountService *AccountService) CreateAccount(account *models.Account) (*models.Account, error) {
	return accountService.accountRepository.CreateAccount(account)
}

func (accountService *AccountService) GetAccountByID(accountId string) (*models.Account, error) {
	return accountService.accountRepository.GetAccountByID(accountId)
}

func (accountService *AccountService) GetAccounts(userId int) ([]*models.Account, error) {
	return accountService.accountRepository.GetAccountsByUserID(userId)
}

func (accountService *AccountService) UpdateAccount(accountId string, accountDTO dto.Account) (*models.Account, error) {
	existingAccount, err := accountService.accountRepository.GetAccountByID(accountId)
	if err != nil {
		return nil, err
	}
	existingAccount.AccountName = accountDTO.AccountName
	existingAccount.Balance = accountDTO.Balance
	existingAccount.Currency = accountDTO.Currency
	updatedAccount, err := accountService.accountRepository.UpdateAccount(existingAccount)
	if err != nil {
		return nil, err
	}
	return updatedAccount, nil
}

func (accountService *AccountService) DeleteAccount(accountID string) (string, error) {
	return accountService.accountRepository.DeleteAccount(accountID)
}
