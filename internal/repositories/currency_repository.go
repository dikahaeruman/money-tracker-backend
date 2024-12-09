package repositories

import (
	"context"
	"database/sql"
	"log"
	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
)

type currencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) interfaces.CurrencyRepositoryInterface {
	return &currencyRepository{db: db}
}

func (r *currencyRepository) GetCurrency(ctx context.Context) ([]*models.Currency, error) {
	query := "SELECT id, currency_code, currency_name FROM currencies"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var currencies []*models.Currency
	for rows.Next() {
		currency := &models.Currency{}
		if err := rows.Scan(&currency.ID, &currency.Code, &currency.Name); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (r *currencyRepository) GetCurrencyByCode(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	query := "SELECT code, name FROM currencies WHERE code = $1"
	err := r.db.QueryRow(query, currency.Code).Scan(&currency.Code, &currency.Name)
	if err != nil {
		return nil, err
	}
	return currency, nil
}
