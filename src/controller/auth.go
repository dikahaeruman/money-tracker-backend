package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	middleware "money-tracker-backend/src/middleware"
	models "money-tracker-backend/src/model"
	util "money-tracker-backend/src/util"
)

// Credentials struct to handle email and password
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles the login request and generates JWT token
func LoginHandler(dbInstance *sql.DB, jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credentials
		var user models.User

		// Parse and validate JSON payload
		if err := c.BindJSON(&creds); err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request payload",
				Data:       map[string]string{"error": err.Error()},
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		// Query to search for user by email
		query := "SELECT id, username, password, email FROM users WHERE email = $1;"
		fmt.Println("Executing query:", query, "with email:", creds.Email)
		// Execute the query
		err := dbInstance.QueryRow(query, creds.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err == sql.ErrNoRows {
			response := util.WriteResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid email or password",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusUnauthorized, response)
			return
		} else if err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid email or password",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		tokenString, err := middleware.CreateJWTToken(jwtKey, creds.Email)
		if err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not generate token",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		refreshTokenString, err := middleware.CreateRefreshToken(jwtKey, creds.Email)
		if err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not generate refresh token",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Set the token as a cookie
		c.SetCookie("token", tokenString, int(middleware.GetDuration().Seconds()), "/", "", false, true)
		c.SetCookie("refresh_token", refreshTokenString, int(middleware.GetRefreshDuration().Seconds()), "/", "", false, true)

		response := util.WriteResponse{
			StatusCode: http.StatusOK,
			Message:    "Login successful",
			Data: map[string]string{
				"token":         tokenString,
				"refresh_token": refreshTokenString,
				"expires_at":    fmt.Sprintf("%v", middleware.GetDuration().Seconds()),
			},
		}

		c.JSON(http.StatusOK, response)

	}
}

func RefreshTokenHandler(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			RefreshToken string `json:"refresh_token"`
		}

		// Bind JSON request body to payload
		if err := c.BindJSON(&payload); err != nil {
			fmt.Printf("check request: %v\n", payload) // Print the raw request body for debugging
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Extract the refresh token from the payload
		refreshTokenStr := payload.RefreshToken
		if refreshTokenStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
			return
		}

		// Parse and validate the refresh token
		token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims"})
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		// Generate a new access token
		newAccessToken, err := middleware.CreateJWTToken(jwtKey, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
			return
		}

		response := util.WriteResponse{
			StatusCode: http.StatusOK,
			Message:    "Refresh Token Generated Successfully",
			Data: map[string]string{
				"token":      newAccessToken,
				"expires_at": fmt.Sprintf("%v", middleware.GetDuration().Seconds()),
			},
		}
		// Set the token as a cookie
		c.SetCookie("token", newAccessToken, int(middleware.GetDuration().Seconds()), "/", "", false, true)

		c.JSON(http.StatusOK, response)
	}
}

// LogoutHandler( handles user logout by clearing the JWT cookie
func LogoutHandler(c *gin.Context) {
	// Clear the token cookie by setting it with an expired value
	c.SetCookie("token", "", -1, "/", "", false, true)

	response := util.WriteResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       []interface{}{}, // Initialize Data as an empty slice
	}

	// Respond to the client indicating the user is logged out
	c.JSON(http.StatusOK, response)

}
