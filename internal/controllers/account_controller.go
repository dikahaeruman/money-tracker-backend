package controllers

import (
	"github.com/gin-gonic/gin"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/services"
	"money-tracker-backend/internal/utils"
	"net/http"
	"strconv"
)

type AccountController struct {
	accountService *services.AccountService
}

func NewAccountController(accountService *services.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

func (accountController *AccountController) GetAllAccounts(c *gin.Context) {
	userID := c.Param("user_id")
	accounts, err := accountController.accountService.GetAllAccounts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get all account"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Account retrieved successfully", accounts))
	return
}

func (accountController *AccountController) GetAccount(c *gin.Context) {
	accountId := c.Param("account_id")
	userIdStr := c.Param("user_id")

	userId, err := strconv.Atoi(userIdStr)

	account, err := accountController.accountService.GetAccount(accountId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Account not found"))
		return
	}
	if account.UserID != userId {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("You don't have permission to access this account"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Account retrieved successfully", account))
	return
}

func (accountController *AccountController) CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	createdAccount, err := accountController.accountService.CreateAccount(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create account"))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Account created successfully", createdAccount))
	return
}
