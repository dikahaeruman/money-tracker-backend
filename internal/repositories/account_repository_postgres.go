package repositories

import (
	"context"
	"database/sql"
	"errors"

	"money-tracker-backend/internal/models"
)

// accountRepositoryPostgres implements AccountRepository interface
type accountRepositoryPostgres struct {
	db *sql.DB
}

// NewAccountRepositoryPostgres creates a new instance of accountRepositoryPostgres
func NewAccountRepositoryPostgres(db *sql.DB) AccountRepository {
	return &accountRepositoryPostgres{db: db}
}

// CreateAccount creates a new account in the database
func (r *accountRepositoryPostgres) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	query := `INSERT INTO accounts (user_id, account_name, balance, currency, created_at) 
                VALUES ($1, $2, $3, $4, NOW()) 
                RETURNING id, user_id, account_name, balance, currency, created_at`
	err := r.db.QueryRowContext(ctx, query, account.UserID, account.AccountName, account.Balance, account.Currency).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountByID retrieves an account by its ID
func (r *accountRepositoryPostgres) GetAccountByID(ctx context.Context, accountID string) (*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency, created_at FROM accounts WHERE id = $1`
	account := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, accountID).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

// GetAccountsByUserID retrieves all accounts for a given user ID
func (r *accountRepositoryPostgres) GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency, created_at FROM accounts WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

// UpdateAccount updates an existing account in the database
func (r *accountRepositoryPostgres) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	query := `UPDATE accounts 
                SET account_name = $1, balance = $2, currency = $3 
                WHERE id = $4 
                RETURNING id, user_id, account_name, balance, currency, created_at`
	err := r.db.QueryRowContext(ctx, query, account.AccountName, account.Balance, account.Currency, account.ID).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

// DeleteAccount removes an account from the database
func (r *accountRepositoryPostgres) DeleteAccount(ctx context.Context, accountID string) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, accountID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}
