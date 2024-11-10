package utils

import (
	"errors"
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

	return duration, nil
}

func CreateJWTToken(email string, userID int) (string, error) {
	// Get the JWT duration and handle potential error
	// expirationDuration, err := GetJWTDuration()
	// if err != nil {
	// 	return "", err // Return an error if getting the duration fails
	// }

	// Use the duration to set the expiration time for the token
	expirationTime := time.Now().Add(200)

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
	// expirationDuration, err := GetRefreshDuration()
	// if err != nil {
	// 	return "", err // Return an error if getting the duration fails
	// }

	// Use the duration to set the expiration time for the token
	expirationTime := time.Now().Add(200)

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
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
