package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"service": "go-order-service",
		"status":  "ok",
	})
}
