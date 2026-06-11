package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

func ListProducts(productRepo *repository.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := productRepo.ListProducts(c.Request.Context())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "获取商品列表失败")
			return
		}
		response.Success(c, products)
	}
}
