package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/cache"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

type CreateOrderRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

func CreateOrder(orderRepo *repository.OrderRepository, idem *cache.IdempotencyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		usernameValue, ok := c.Get("username")
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID, ok := userIDValue.(int64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		username, ok := usernameValue.(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "idempotency key is required")
			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request")
			return
		}

		if req.ProductID <= 0 {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "product_id must be greater than 0")
			return
		}

		if req.Quantity <= 0 {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "quantity must be greater than 0")
			return
		}

		acquired, err := idem.Acquire(c.Request.Context(), userID, idempotencyKey)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "idempotency service unavailable")
			return
		}

		if !acquired {
			response.Error(c, http.StatusConflict, http.StatusConflict, "duplicate request")
			return
		}

		order, err := orderRepo.CreateOrder(c.Request.Context(), userID, username, req.ProductID, req.Quantity)
		if err != nil {
			_ = idem.Release(c.Request.Context(), userID, idempotencyKey)

			switch {
			case errors.Is(err, repository.ErrProductNotFound):
				response.Error(c, http.StatusNotFound, http.StatusNotFound, "product not found")
			case errors.Is(err, repository.ErrInsufficientStock):
				response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "insufficient stock")
			default:
				response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		_ = idem.MarkSuccess(c.Request.Context(), userID, idempotencyKey, order.OrderNo)

		response.Success(c, order)
	}
}
