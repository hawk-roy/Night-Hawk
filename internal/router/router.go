package router

import (
	"github.com/gin-gonic/gin"
	"github.com/teng-lei-tfs/Night-Hawk/internal/handler"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)
	}

	return r
}
