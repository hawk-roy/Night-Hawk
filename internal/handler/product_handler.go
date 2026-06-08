package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
)

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
