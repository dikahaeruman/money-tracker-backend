package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
)

type accountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new instance of accountRepository
func NewAccountRepository(db *sql.DB) interfaces.AccountRepositoryInterface {
	return &accountRepository{db: db}
}

// CreateAccount creates a new account in the database
func (r *accountRepository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	existingAccount, err := r.GetAccountByName(ctx, account.AccountName)
	if err != nil && err != sql.ErrNoRows {
		// If the error is due to a database error, handle it
		fmt.Println("error on get: ", err)
		return nil, err
	}

	if existingAccount != nil {
		return nil, fmt.Errorf("account with name %s already exists", account.AccountName)
	}

	query := `INSERT INTO accounts (user_id, account_name, balance, currency_id, created_at) 
                VALUES ($1, $2, $3, $4, NOW()) 
                RETURNING id, user_id, account_name, balance, currency_id, created_at`

	err = r.db.QueryRowContext(ctx, query, account.UserID, account.AccountName, account.Balance, account.CurrencyID).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.CurrencyID, &account.CreatedAt)
	if err != nil {
		fmt.Println("error on insert: ", err)
		return nil, err
	}
	return account, nil
}

// GetAccountByID retrieves an account by its ID
func (r *accountRepository) GetAccountByID(ctx context.Context, accountID string) (*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency_id, created_at FROM accounts WHERE id = $1`
	account := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, accountID).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.CurrencyID, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

// GetAccountsByUserID retrieves all accounts for a given user ID
func (r *accountRepository) GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error) {
	fmt.Println("success get accounts id!: ", userID)

	query := `SELECT 
					accounts.id, 
					accounts.user_id, 
					accounts.account_name, 
					accounts.balance, 
					accounts.currency_id,
					currencies.currency_code,
					accounts.created_at,
					accounts.updated_at
				FROM
					accounts
				JOIN currencies on accounts.currency_id = currencies.id
				WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		fmt.Println("error get accounts id!: ", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		currency := &models.Currency{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountName,
			&account.Balance,
			&account.CurrencyID,
			&currency.Code,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			fmt.Println("error get accounts!: ", err)
			return nil, err
		}
		account.CurrencyCode = currency.Code
		fmt.Println("success get accounts!: ", account)
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("error get accountssssss!: ", err)
		return nil, err
	}

	return accounts, nil
}

// GetAccountByName retrieves an account by its name
func (r *accountRepository) GetAccountByName(ctx context.Context, accountName string) (*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency_id, created_at FROM accounts WHERE account_name = $1`
	account := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, accountName).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.CurrencyID, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// UpdateAccount updates an existing account in the database
func (r *accountRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	query := `UPDATE accounts 
                SET account_name = $1, balance = $2, currency_id = $3 
                WHERE id = $4 
                RETURNING id, user_id, account_name, balance, currency_id, created_at`
	err := r.db.QueryRowContext(ctx, query, account.AccountName, account.Balance, account.CurrencyID, account.ID).
		Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.CurrencyID, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

// DeleteAccount removes an account from the database
func (r *accountRepository) DeleteAccount(ctx context.Context, accountID string) error {
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
