package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/handler"
	"github.com/hawk-roy/Night-Hawk/internal/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)

		users := api.Group("/users")
		{
			users.POST("/register", handler.RegisterUser)
			users.POST("/login", handler.Login)
		}

		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/users/me", handler.Me)
		}
	}

	return r
}
