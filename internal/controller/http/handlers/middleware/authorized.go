package middleware

import (
	"net/http"
	"strings"

	"experiment.io/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authorized(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": entity.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": entity.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
