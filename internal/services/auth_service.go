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

func (s *Service) Authenticate(email, password string) (string, error) {
	user, err := s.userRepo.FindPasswordByEmail(email)
	if err != nil {
		return "", err
	}

	match, err := utils.CheckPasswordHash([]byte(user.Password), password)
	if err != nil {
		log.Printf("Error checking password: %v", err)
		return "", err
	}
	if !match {
		log.Printf("Invalid password for user: %s", email)
		return "", errors.New("invalid email or password")
	}

	return utils.CreateJWTToken(user.Email)
}

func (s *Service) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	return utils.CreateJWTToken(claims.Email)
}

func (s *Service) VerifyToken(token string) (*utils.Claims, error) {
	return utils.VerifyToken(token)
}
