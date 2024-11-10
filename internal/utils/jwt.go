package utils

import (
	"errors"
	"fmt"
	"log"

	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email  string `json:"email"`
	UserID int    `json:"id"`
	jwt.StandardClaims
}

func GetJWTKey() string {
	return os.Getenv("JWT_SECRET")
}

func GetJWTDuration() (time.Duration, error) {
	// Get the value from the environment variable
	strDuration := os.Getenv("LOGIN_EXPIRATION_DURATION")

	if strDuration == "" {
		return 0, errors.New("LOGIN_EXPIRATION_DURATION not set")
	}

	// Parse the duration string (e.g., "24h" or "30m")
	duration, err := time.ParseDuration(strDuration)
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}
	log.Printf("JWT Duration: %s", duration.String())

	return duration, nil
}

func GetRefreshDuration() (time.Duration, error) {
	// Get the value from the environment variable
	strDuration := os.Getenv("REFRESH_EXPIRATION_DURATION")

	if strDuration == "" {
		return 0, errors.New("REFRESH_EXPIRATION_DURATION not set")
	}

	// Parse the duration string (e.g., "24h" or "30m")
	duration, err := time.ParseDuration(strDuration)
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}
	log.Printf("Refresh Duration: %s", duration.String())
	return duration, nil
}

func CreateJWTToken(email string, userID int) (string, error) {
	// Get the JWT duration and handle potential error
	expirationDuration, err := GetJWTDuration()
	if err != nil {
		return "", err // Return an error if getting the duration fails
	}

	// Use the duration to set the expiration time for the token
	expirationTime := time.Now().Add(time.Duration(int64(expirationDuration.Seconds())) * time.Second)
	log.Printf("Expiration Time: %s", expirationTime.String())

	claims := &Claims{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(GetJWTKey()))
	if err != nil {
		return "", err // Return any error encountered during signing
	}

	return signedToken, nil
}

func CreateRefreshToken(email string, userID int) (string, error) {
	// Get the Refresh duration and handle potential error
	expirationDuration, err := GetRefreshDuration()
	if err != nil {
		return "", err // Return an error if getting the duration fails
	}

	// Use the duration to set the expiration time for the token
	expirationTime := time.Now().Add(time.Duration(int64(expirationDuration.Seconds())) * time.Second)
	log.Printf("Expiration Time: %s", expirationTime.String())

	claims := &Claims{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(GetJWTKey()))
	if err != nil {
		return "", err // Return any error encountered during signing
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTKey()), nil
	})
	log.Printf("Current Unix Time: %d", time.Now().Unix())
	log.Printf("Token Expiry: %d", token.Claims.(*Claims).ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Optionally, check specific claims for refresh token (if you have such requirements)
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
