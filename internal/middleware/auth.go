package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/auth"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()

		//收尾工作

	}
}
