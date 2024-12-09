package interfaces

import (
	"context"
	"money-tracker-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CurrencyRepositoryInterface interface {
	GetCurrency(ctx context.Context) ([]*models.Currency, error)
	GetCurrencyByCode(ctx context.Context, currency *models.Currency) (*models.Currency, error)
}

type CurrencyServiceInterface interface {
	GetCurrency(ctx context.Context) ([]*models.Currency, error)
	GetCurrencyByCode(ctx context.Context, currency *models.Currency) (*models.Currency, error)
}

type CurrencyControllerInterface interface {
	GetCurrency(ctx *gin.Context)
	GetCurrencyByCode(ctx *gin.Context)
}
