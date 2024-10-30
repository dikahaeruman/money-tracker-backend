package repositories

import (
	"database/sql"
	"money-tracker-backend/internal/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetAll(userId string) ([]*models.Account, error) {
	query := `
        SELECT id, user_id, account_name, balance, currency, created_at
        FROM accounts WHERE user_id = $1`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	accounts := make([]*models.Account, 0)
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *AccountRepository) Get(accountId string) (*models.Account, error) {
	query := `
        SELECT id, user_id, account_name, balance, currency, created_at
        FROM accounts WHERE id = $1`
	account := &models.Account{}
	err := r.db.QueryRow(query, accountId).Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) Create(account *models.Account) (*models.Account, error) {
	query := `
        INSERT INTO accounts (user_id, account_name, balance, currency, created_at)
    VALUES ($1, $2, $3, $4, NOW())
    RETURNING id, user_id, account_name, balance, currency, created_at`
	err := r.db.QueryRow(query, account.UserID, account.AccountName, account.Balance, account.Currency).Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}
