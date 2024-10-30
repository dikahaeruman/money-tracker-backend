package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"money-tracker-backend/internal/utils"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Missing token"))
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token"))
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
