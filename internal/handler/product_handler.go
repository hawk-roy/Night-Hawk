package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/model"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
)

var orderProducts = []model.Product{
	{
		ID:          1,
		Name:        "Go Backend Course",
		Description: "A practical Go backend course",
		Price:       19900,
		Stock:       100,
	},
	{
		ID:          2,
		Name:        "API Design Handbook",
		Description: "A handbook for designing backend APIs",
		Price:       9900,
		Stock:       50,
	},
	{
		ID:          3,
		Name:        "Cloud Native Starter Kit",
		Description: "A starter kit for cloud native applications",
		Price:       29900,
		Stock:       30,
	},
}

func GetProductByID(productID int64) (*model.Product, bool) {
	for i := range orderProducts {
		if orderProducts[i].ID == productID {
			return &orderProducts[i], true
		}
	}

	return nil, false
}

func ListProducts(productRepo *repository.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := productRepo.ListProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "failed to list products",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    products,
		})
	}
}
