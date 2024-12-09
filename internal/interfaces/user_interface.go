// internal/interfaces/user_service.go

package interfaces

import (
	"money-tracker-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(email string) (*models.User, error)
	SearchByEmail(username string) (*models.User, error)
}

type UserControllerInterface interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
}
