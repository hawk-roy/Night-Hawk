package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/cache"
	"github.com/hawk-roy/Night-Hawk/internal/handler"
	"github.com/hawk-roy/Night-Hawk/internal/middleware"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"github.com/redis/go-redis/v9"
)

func NewRouter(mysqlDB *sql.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	userRepo := repository.NewUserRepository(mysqlDB)
	productRepo := repository.NewProductRepository(mysqlDB)
	orderRepo := repository.NewOrderRepository(mysqlDB)
	idempotencyManager := cache.NewIdempotencyManager(redisClient)

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)
		api.GET("/health/db", handler.DBHealthCheck(mysqlDB))
		api.GET("/health/redis", handler.RedisHealthCheck(redisClient))
		api.GET("/products", handler.ListProducts(productRepo))

		users := api.Group("/users")
		{
			users.POST("/register", handler.Register(userRepo))
			users.POST("/login", handler.Login(userRepo))
		}

		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/users/me", handler.Me)
			authGroup.POST("/orders", handler.CreateOrder(orderRepo, idempotencyManager))
		}
	}

	return r
}
