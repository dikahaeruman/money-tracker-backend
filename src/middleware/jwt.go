package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims struct to handle JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GetDuration() time.Duration {
	durationStr := os.Getenv("LOGIN_EXPIRATION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf("Invalid duration: %v", err)
	}
	return duration
}

func GetRefreshDuration() time.Duration {
	durationStr := os.Getenv("REFRESH_EXPIRATION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf("Invalid duration: %v", err)
	}
	return duration
}

// JWTMiddleware validates JWT tokens from cookies
func JWTMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set claims into gin context for further use
		c.Set("claims", claims)

		c.Next()
	}
}

func CreateJWTToken(jwtKey []byte, email string) (string, error) {
	expirationTime := time.Now().Add(GetDuration())
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(jwtKey []byte, email string) (string, error) {
	expirationTime := time.Now().Add(GetRefreshDuration()).Unix()
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["email"] = email
	rtClaims["exp"] = expirationTime

	rt, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return rt, nil
}
