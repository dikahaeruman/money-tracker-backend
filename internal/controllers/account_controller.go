package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-tracker-backend/internal/dto"
	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/utils"
)

// AccountController handles HTTP requests related to accounts
type AccountController struct {
	accountService interfaces.AccountServiceInterface
}

// NewAccountController creates a new instance of AccountController
func NewAccountController(accountService interfaces.AccountServiceInterface) interfaces.AccountControllerInterface {
	return &AccountController{accountService: accountService}
}

// CreateAccount handles the creation of a new account
func (ac *AccountController) CreateAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	account.UserID = userID.(int)

	if err := account.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	createdAccount, err := ac.accountService.CreateAccount(c.Request.Context(), &account)
	if err != nil {
		fmt.Println("error on create controller: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create account"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Account created successfully", createdAccount))
}

// GetAccountByID retrieves an account by its ID
func (ac *AccountController) GetAccountByID(c *gin.Context) {
	accountID := c.Param("account_id")

	account, err := ac.accountService.GetAccountByID(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Account not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Account retrieved successfully", account))
}

// GetAccounts retrieves all accounts for the authenticated user
func (ac *AccountController) GetAccounts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	accounts, err := ac.accountService.GetAccounts(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get accounts"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Accounts retrieved successfully", accounts))
}

// UpdateAccount updates an existing account
func (ac *AccountController) UpdateAccount(c *gin.Context) {
	accountID := c.Param("account_id")

	var accountDTO dto.Account
	if err := c.ShouldBindJSON(&accountDTO); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	updatedAccount, err := ac.accountService.UpdateAccount(c.Request.Context(), accountID, accountDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update account"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Account updated successfully", updatedAccount))
}

// DeleteAccount deletes an account by its ID
func (ac *AccountController) DeleteAccount(c *gin.Context) {
	accountID := c.Param("account_id")

	err := ac.accountService.DeleteAccount(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete account"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Account deleted successfully", nil))
}
