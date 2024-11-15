package services

import (
	"errors"
	"log"

	_ "money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
	"money-tracker-backend/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Authenticate(email, password string) (string, string, error) { // Updated return type
	user, err := s.userRepo.FindPasswordByEmail(email)
	if err != nil {
		return "", "", err
	}

	match, err := utils.CheckPasswordHash([]byte(user.Password), password)
	if err != nil {
		log.Printf("Error checking password: %v", err)
		return "", "", err
	}
	if !match {
		log.Printf("Invalid password for user: %s", email)
		return "", "", errors.New("invalid email or password")
	}

	jwtToken, err := utils.CreateJWTToken(user.Email, user.ID)
	if err != nil {
		log.Printf("Error creating JWT token: %v", err)
		return "", "", err
	}
	refreshToken, err := utils.CreateRefreshToken(user.Email, user.ID) // Ensure this line is correct
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		return "", "", err
	}

	return jwtToken, refreshToken, nil // Updated return statement
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}
	log.Printf("Refresh token: %s", refreshToken)
	log.Printf("Claims: %v", claims)
	return utils.CreateJWTToken(claims.Email, claims.UserID)
}

func (s *AuthService) VerifyToken(token string) (*utils.Claims, error) {
	return utils.VerifyToken(token)
}
