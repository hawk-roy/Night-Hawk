package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/model"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

type MockPayOrderRequest struct {
	OrderID int64  `json:"order_id"`
	Result  string `json:"result"`
}

type MockPayOrderResponse struct {
	Payment     *model.Payment `json:"payment"`
	OrderStatus string         `json:"order_status"`
}

func MockPayOrder(paymentRepo *repository.PaymentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID, ok := userIDValue.(int64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
			return
		}

		var orderReq MockPayOrderRequest

		if err := c.ShouldBindJSON(&orderReq); err != nil {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request")
			return
		}

		if orderReq.OrderID <= 0 {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "order_id must be greater than 0")
			return
		}

		if orderReq.Result != model.PaymentStatusSuccess && orderReq.Result != model.PaymentStatusFailed {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid payment result")
			return
		}

		payment, orderStatus, err := paymentRepo.MockPayOrder(c.Request.Context(), userID, orderReq.OrderID, orderReq.Result)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrOrderNotFound):
				response.Error(c, http.StatusNotFound, http.StatusNotFound, "order not found")
			case errors.Is(err, repository.ErrOrderNotPendingPayment):
				response.Error(c, http.StatusConflict, http.StatusConflict, "order is not pending payment")
			case errors.Is(err, repository.ErrInvalidPaymentResult):
				response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid payment result")
			default:
				response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		response.Success(c, MockPayOrderResponse{
			Payment:     payment,
			OrderStatus: orderStatus,
		})

	}
}
