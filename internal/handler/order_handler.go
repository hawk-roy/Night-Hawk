package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/model"
)

var (
	orderIDCounter int64
	orders         []model.Order
	orderMu        sync.Mutex
)

type CreateOrderRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

func CreateOrder(c *gin.Context) {
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

	product, ok := GetProductByID(req.ProductID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "product not found",
			"data":    nil,
		})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "insufficient stock",
			"data":    nil,
		})
		return
	}

	orderMu.Lock()
	defer orderMu.Unlock()

	orderIDCounter++

	order := model.Order{
		ID:          orderIDCounter,
		UserID:      userID,
		Username:    username,
		ProductID:   product.ID,
		ProductName: product.Name,
		UnitPrice:   product.Price,
		Quantity:    req.Quantity,
		TotalAmount: product.Price * req.Quantity,
		Status:      model.OrderStatusPendingPayment,
		CreatedAt:   time.Now(),
	}

	orders = append(orders, order)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}
