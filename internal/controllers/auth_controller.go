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

	// Authenticate user and get tokens
	token, refreshToken, err := ac.service.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid email or password"))
		return
	}

	log.Printf("Token: %s, Refresh token: %s", token, refreshToken)

	// Get durations and handle errors
	jwtDuration, err := utils.GetJWTDuration()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not get token duration"))
		return
	}

	refreshDuration, err := utils.GetRefreshDuration()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not get refresh token duration"))
		return
	}

	log.Printf("Token: %s, Refresh token: %s", token, refreshToken)

	// Set cookies with correct expiration
	c.SetCookie("token", token, int(jwtDuration.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, int(refreshDuration.Seconds()), "/", "", true, true)

	// Respond with success
	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
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

	jwtDuration, err := utils.GetJWTDuration()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not get token duration"))
		return
	}

	refreshDuration, err := utils.GetRefreshDuration()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not get refresh token duration"))
		return
	}

	c.SetCookie("token", newToken, int(refreshDuration.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed successfully", gin.H{
		"token":      newToken,
		"expires_at": jwtDuration.Seconds(),
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
	c.Header("Clear-Site-Data", "\"cookies\"")
	c.SetSameSite(http.SameSiteStrictMode)
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.JSON(http.StatusOK, utils.SuccessResponse("Logout successful", nil))
}
