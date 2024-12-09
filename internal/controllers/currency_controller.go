package controllers

import (
	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	currencyService interfaces.CurrencyServiceInterface
}

func NewCurrencyController(currencyService interfaces.CurrencyServiceInterface) interfaces.CurrencyControllerInterface {
	return &CurrencyController{currencyService: currencyService}
}

func (c *CurrencyController) GetCurrency(ctx *gin.Context) {
	currencies, err := c.currencyService.GetCurrency(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get currencies"))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Currencies retrieved successfully", currencies))
}

func (c *CurrencyController) GetCurrencyByCode(ctx *gin.Context) {
	c.currencyService.GetCurrencyByCode(ctx, &models.Currency{Code: ctx.Query("code")})
}
