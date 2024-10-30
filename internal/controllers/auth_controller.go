package controllers

import (
	"log"
	"money-tracker-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/utils"
)

type AuthController struct {
	service *services.Service
}

func NewAuthController(service *services.Service) *AuthController {
	return &AuthController{service: service}
}

func (ac *AuthController) Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	token, err := ac.service.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid email or password"))
		return
	}

	c.SetCookie("token", token, int(utils.GetJWTDuration().Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", gin.H{
		"token":      token,
		"expires_at": utils.GetJWTDuration().Seconds(),
	}))
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
		return
	}

	newToken, err := ac.service.RefreshToken(payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid refresh token"))
		return
	}

	c.SetCookie("token", newToken, int(utils.GetJWTDuration().Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed successfully", gin.H{
		"token":      newToken,
		"expires_at": utils.GetJWTDuration().Seconds(),
	}))
}

func (ac *AuthController) VerifyToken(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Missing token"))
		return
	}

	claims, err := ac.service.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Token verified", claims))
}

func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, utils.SuccessResponse("Logout successful", nil))
}
