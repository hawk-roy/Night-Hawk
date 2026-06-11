package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/response"
	"github.com/redis/go-redis/v9"
)

func RedisHealthCheck(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "redis unavailable")
			return
		}

		response.Success(c, gin.H{
			"cache":  "redis",
			"status": "ok",
		})

	}
}
