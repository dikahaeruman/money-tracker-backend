package controllers

import (
	"github.com/gin-gonic/gin"
	"money-tracker-backend/internal/dto"
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

func (accountController *AccountController) GetAccountByID(c *gin.Context) {
	accountId := c.Param("account_id")

	account, err := accountController.accountService.GetAccountByID(accountId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Account not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Account retrieved successfully", account))
	return
}

func (accountController *AccountController) GetAccountsByUserID(c *gin.Context) {
	userIDString := c.Param("user_id")
	userID, err := strconv.Atoi(userIDString)
	accounts, err := accountController.accountService.GetAccountsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get all account"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Account retrieved successfully", accounts))
	return
}

func (accountController *AccountController) UpdateAccount(c *gin.Context) {
	accountId := c.Param("account_id")
	var accountDTO dto.Account

	if err := c.ShouldBindJSON(&accountDTO); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}
	account, err := accountController.accountService.GetAccountByID(accountId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Account not found"))
		return
	}

	account.ID = accountId
	account.AccountName = accountDTO.AccountName
	account.Balance = accountDTO.Balance
	account.Currency = accountDTO.Currency

	updatedAccount, err := accountController.accountService.UpdateAccount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update account"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Account updated successfully", updatedAccount))
	return
}

func (accountController *AccountController) DeleteAccount(c *gin.Context) {
	accountId := c.Param("account_id")

	deletedId, err := accountController.accountService.DeleteAccount(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete account"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Account deleted successfully", deletedId))
	return
}
