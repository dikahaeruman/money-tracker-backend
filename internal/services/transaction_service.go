package services

import (
	"context"
	"errors"
	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transactionDTO *dto.Transaction) (*models.Transaction, error)
}

type transactionService struct {
	repo           repositories.TransactionRepository
	accountService AccountService
}

func NewTransactionService(repo repositories.TransactionRepository, accountService AccountService) TransactionService {
	return &transactionService{
		repo:           repo,
		accountService: accountService,
	}
}

func (transactionService *transactionService) CreateTransaction(ctx context.Context, transactionDTO *dto.Transaction) (*models.Transaction, error) {
	var transaction models.Transaction
	account, err := transactionService.accountService.GetAccountByID(ctx, transactionDTO.AccountID)
	if err != nil {
		return nil, err
	}
	transaction.AccountID = transactionDTO.AccountID
	transaction.TransactionType = transactionDTO.TransactionType
	transaction.BalanceBefore = account.Balance
	transaction.Amount = transactionDTO.Amount
	transaction.Description = transactionDTO.Description
	transaction.TransactionDate = transactionDTO.TransactionDate
	if transaction.TransactionType == "credit" {
		account.Balance += transaction.Amount
	}
	if transaction.TransactionType == "debit" {
		if transaction.Amount > account.Balance {
			return nil, errors.New("insufficient balance")
		}
		account.Balance -= transaction.Amount
	}

	account.Balance = transaction.BalanceAfter
	_, err = transactionService.accountService.UpdateAccount(ctx, account.ID, dto.Account{
		AccountName: account.AccountName,
		Balance:     account.Balance,
		Currency:    account.Currency,
	})
	if err != nil {
		return nil, err
	}

	return transactionService.repo.CreateTransaction(ctx, &transaction)
}
