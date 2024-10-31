package repositories

import (
	"database/sql"
	"money-tracker-backend/internal/models"
)

type accountRepositoryPostgres struct {
	db *sql.DB
}

func NewAccountRepositoryPostgres(db *sql.DB) AccountRepository {
	return &accountRepositoryPostgres{db: db}
}

func (r *accountRepositoryPostgres) CreateAccount(account *models.Account) (*models.Account, error) {
	query := `INSERT INTO accounts (user_id, account_name, balance, currency, created_at) 
				VALUES ($1, $2, $3, $4, NOW()) 
				RETURNING id, user_id, account_name, balance, currency, created_at`
	err := r.db.QueryRow(query, account.UserID, account.AccountName, account.Balance, account.Currency).Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepositoryPostgres) GetAccountByID(accountID string) (*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency, created_at FROM accounts WHERE id = $1;`
	account := &models.Account{}
	err := r.db.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.AccountName, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepositoryPostgres) GetAccountsByUserID(userID int) ([]*models.Account, error) {
	query := `SELECT id, user_id, account_name, balance, currency, created_at FROM accounts WHERE user_id = $1;`
	rows, err := r.db.Query(query, userID)
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

func (r *accountRepositoryPostgres) UpdateAccount(account *models.Account) (*models.Account, error) {
	query := `UPDATE accounts 
				SET account_name = $1, balance = $2, currency = $3 
				WHERE id = $4 
				RETURNING id, user_id, account_name, balance, currency, created_at;`
	updatedAccount := &models.Account{}
	err := r.db.QueryRow(query, account.AccountName, account.Balance, account.Currency, account.ID).Scan(
		&updatedAccount.ID,
		&updatedAccount.UserID,
		&updatedAccount.AccountName,
		&updatedAccount.Balance,
		&updatedAccount.Currency,
		&updatedAccount.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepositoryPostgres) DeleteAccount(accountID string) (string, error) {
	query := `DELETE FROM accounts WHERE id = $1 RETURNING id;`
	var deletedAccountID string
	err := r.db.QueryRow(query, accountID).Scan(&deletedAccountID)
	if err != nil {
		return `Error deleting accounts`, err
	}
	return deletedAccountID, nil
}
