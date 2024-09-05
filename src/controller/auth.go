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

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(dbInstance *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credentials
		var user models.User

		jwtKey, ok := util.GetStringFromContext(c, "jwtKey")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT key not found or invalid type"})
			return
		}

		if err := c.BindJSON(&creds); err != nil {
			response := util.WriteResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request payload",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		query := "SELECT id, username, password, email FROM users WHERE email = $1;"
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

		c.SetCookie("token", tokenString, int(middleware.GetDuration().Seconds()), "/", "", false, true)

		response := util.WriteResponse{
			StatusCode: http.StatusOK,
			Message:    "Login successful",
			Data: map[string]string{
				"token":      tokenString,
				"expires_at": fmt.Sprintf("%v", middleware.GetDuration().Seconds()),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func RefreshTokenHandler(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}

	jwtKey, ok := util.GetStringFromContext(c, "jwtKey")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT key not found or invalid type"})
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		fmt.Printf("check request: %v\n", payload)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	refreshTokenStr := payload.RefreshToken
	if refreshTokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtKey), nil
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email in token"})
		return
	}

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
	c.SetCookie("token", newAccessToken, int(middleware.GetDuration().Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

func VerifyTokenHandler(jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		claims, err := middleware.VerifyToken(cookie, jwtKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"claims": claims})
	}
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)

	response := util.WriteResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       []interface{}{},
	}

	c.JSON(http.StatusOK, response)
}
