package utils

import (
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

func GetJWTDuration() time.Duration {
	return 15 * time.Minute
}

func CreateJWTToken(email string, userID int) (string, error) {
	expirationTime := time.Now().Add(GetJWTDuration())
	claims := &Claims{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetJWTKey()))
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
