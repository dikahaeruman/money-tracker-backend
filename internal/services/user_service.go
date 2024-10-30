package services

import (
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/repositories"
	"money-tracker-backend/internal/utils"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (userService *UserService) CreateUser(user *models.User) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return userService.userRepo.Create(user)
}

func (userService *UserService) GetAllUsers() ([]models.User, error) {
	return userService.userRepo.FindAll()
}

func (userService *UserService) SearchUserByUsername(username string) (*models.User, error) {
	return userService.userRepo.FindByUsername(username)
}
