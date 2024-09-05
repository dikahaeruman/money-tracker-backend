package middleware

import (
	"fmt"
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

func GetJWTKey() string {
	key := os.Getenv("JWT_KEY")
	log.Printf("JWT_KEY: %s", key)
	return key
}

func GetDuration() time.Duration {
	durationStr := os.Getenv("LOGIN_EXPIRATION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf("Invalid duration: %v", err)
	}
	return duration
}

func GetRefreshDuration() time.Duration {
	durationStr := os.Getenv("REFRESH_EXPIRATION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf("Invalid duration: %v", err)
	}
	return duration
}

func JWTMiddleware(jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			log.Println("Missing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		claims, err := VerifyToken(cookie, jwtKey)

		log.Println("claims handler:", claims)

		if err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func VerifyToken(tokenString string, jwtKey string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	// Log the token for debugging
	log.Printf("Token String: %v", tokenString)
	log.Printf("Parsed Token: %v", token)

	// Check if parsing failed
	if err != nil {
		// Check for expiration error
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			log.Println("Token expired")
			return nil, fmt.Errorf("token expired")
		}
		// Any other error
		log.Println("Token validation error:", err)
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		log.Println("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	// Return the claims from the token
	return claims, nil
}

func CreateJWTToken(jwtKey string, email string) (string, error) {
	expirationTime := time.Now().Add(GetDuration())
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(jwtKey string, email string) (string, error) {
	expirationTime := time.Now().Add(GetRefreshDuration())
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token using the claims
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the key
	rt, err := refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return rt, nil
}
