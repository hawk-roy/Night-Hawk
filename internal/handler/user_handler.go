package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/model"
)

var (
	userIDCounter int
	users         []model.User
	userMu        sync.Mutex
)

func RegisterUser(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request",
		})
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	userIDCounter++

	user := model.User{
		ID:       userIDCounter,
		Username: req.Username,
	}

	users = append(users, user)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    user,
	})
}
