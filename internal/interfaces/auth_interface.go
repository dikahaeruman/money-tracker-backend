// internal/interfaces/auth_service_interface.go

package interfaces

import (
	"money-tracker-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthService defines the interface for authentication-related operations
type AuthServiceInterface interface {
	Authenticate(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	VerifyToken(token string) (*utils.Claims, error)
}

type AuthControllerInterface interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	VerifyToken(c *gin.Context)
	Logout(c *gin.Context)
}
