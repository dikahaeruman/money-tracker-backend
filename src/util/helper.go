package util

import "github.com/gin-gonic/gin"

func GetStringFromContext(c *gin.Context, key string) (string, bool) {
	value, exists := c.Get(key)
	if !exists {
		return "", false
	}

	strValue, ok := value.(string)
	return strValue, ok
}
