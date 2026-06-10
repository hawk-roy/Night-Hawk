package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/cache"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
)

type CreateOrderRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

func CreateOrder(orderRepo *repository.OrderRepository, idem *cache.IdempotencyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			return
		}

		usernameValue, ok := c.Get("username")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			return
		}

		userID, ok := userIDValue.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			return
		}

		username, ok := usernameValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			return
		}

		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "idempotency key is required",
				"data":    nil,
			})
			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "invalid request",
				"data":    nil,
			})
			return
		}

		if req.ProductID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "product_id must be greater than 0",
				"data":    nil,
			})
			return
		}

		if req.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "quantity must be greater than 0",
				"data":    nil,
			})
			return
		}

		acquired, err := idem.Acquire(c.Request.Context(), userID, idempotencyKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "idempotency service unavailable",
				"data":    nil,
			})
			return
		}

		if !acquired {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "duplicate request",
				"data":    nil,
			})
			return
		}

		order, err := orderRepo.CreateOrder(c.Request.Context(), userID, username, req.ProductID, req.Quantity)
		if err != nil {
			_ = idem.Release(c.Request.Context(), userID, idempotencyKey)

			switch {
			case errors.Is(err, repository.ErrProductNotFound):
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "product not found",
					"data":    nil,
				})
			case errors.Is(err, repository.ErrInsufficientStock):
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "insufficient stock",
					"data":    nil,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "internal server error",
					"data":    nil,
				})
			}
			return
		}

		_ = idem.MarkSuccess(c.Request.Context(), userID, idempotencyKey, order.OrderNo)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    order,
		})
	}
}
