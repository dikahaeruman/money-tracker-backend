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
			response := util.WriteResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not generate token",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Set the token as a cookie
		c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "", false, true)

		response := util.WriteResponse{
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

	response := util.WriteResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       []interface{}{}, // Initialize Data as an empty slice
	}

	// Respond to the client indicating the user is logged out
	c.JSON(http.StatusOK, response)

}
