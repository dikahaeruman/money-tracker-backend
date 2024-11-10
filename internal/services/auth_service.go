package services

import (
	"errors"
	"log"

	_ "money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
	"money-tracker-backend/internal/utils"
)

type Service struct {
	userRepo *repositories.UserRepository
}

func NewService(userRepo *repositories.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Authenticate(email, password string) (string, string, error) { // Updated return type
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

func (s *Service) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}
	log.Printf("Refresh token: %s", refreshToken)
	log.Printf("Claims: %v", claims)
	return utils.CreateJWTToken(claims.Email, claims.UserID)
}

func (s *Service) VerifyToken(token string) (*utils.Claims, error) {
	return utils.VerifyToken(token)
}
