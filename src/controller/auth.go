package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	middleware "money-tracker-backend/src/middleware"
	models "money-tracker-backend/src/model"
)

// Credentials struct to handle email and password
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Todo make this helper later
type WriteResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// LoginHandler handles the login request and generates JWT token
func LoginHandler(dbInstance *sql.DB, jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credentials
		var user models.User

		// Parse and validate JSON payload
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Query to search for user by email
		query := "SELECT id, username, password, email FROM users WHERE email = $1;"
		fmt.Println("Executing query:", query, "with email:", creds.Email)
		// Execute the query
		err := dbInstance.QueryRow(query, creds.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Create the JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &middleware.Claims{
			Email: creds.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Set the token as a cookie
		c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "", false, true)

		response := WriteResponse{
			StatusCode: http.StatusOK,
			Message:    "Login successful",
			Data:       map[string]string{"token": tokenString},
		}

		c.JSON(http.StatusOK, response)

	}
}

// LogoutHandler( handles user logout by clearing the JWT cookie
func LogoutHandler(c *gin.Context) {
	// Clear the token cookie by setting it with an expired value
	c.SetCookie("token", "", -1, "/", "", false, true)

	// Respond to the client indicating the user is logged out
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
