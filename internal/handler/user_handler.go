package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/auth"
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
		Password: req.Password,
	}

	users = append(users, user)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request",
		})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "username is required",
		})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "password is required",
		})
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	for _, user := range users {
		if user.Username == req.Username {
			if user.Password != req.Password {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "invalid username or password",
				})
				return
			}

			token, err := auth.GenerateToken(int64(user.ID), user.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "failed to generate token",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"token": token,
				},
			})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    401,
		"message": "invalid username or password",
	})
}
