package services

import (
	"context"
	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
)

type CurrencyService struct {
	repo interfaces.CurrencyRepositoryInterface
}

func NewCurrencyService(repo interfaces.CurrencyRepositoryInterface) interfaces.CurrencyServiceInterface {
	return &CurrencyService{repo: repo}
}

func (s *CurrencyService) GetCurrency(ctx context.Context) ([]*models.Currency, error) {
	return s.repo.GetCurrency(ctx)
}

func (s *CurrencyService) GetCurrencyByCode(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	return s.repo.GetCurrencyByCode(ctx, currency)
}
