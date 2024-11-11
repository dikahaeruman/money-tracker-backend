package repositories

import (
	"context"
	"database/sql"

	"money-tracker-backend/internal/models"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
}

type TransactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepositoryImpl(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (repo *TransactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	query := `INSERT INTO transactions (account_id, transaction_type, amount, balance_before, balance_after, description, transaction_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING *`
	err := repo.db.QueryRowContext(ctx, query, transaction.AccountID, transaction.TransactionType, transaction.Amount, transaction.BalanceBefore, transaction.BalanceAfter, transaction.Description, transaction.TransactionDate).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.BalanceBefore,
		&transaction.BalanceAfter,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
