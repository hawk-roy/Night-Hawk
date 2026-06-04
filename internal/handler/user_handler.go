package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/auth"
	"github.com/hawk-roy/Night-Hawk/internal/model"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	userIDCounter int
	users         []model.User
	userMu        sync.Mutex
)

func Register(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RegisterRequest

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

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "failed to hash password",
			})
			return
		}

		user, err := userRepo.CreateUser(c.Request.Context(), req.Username, string(passwordHash))
		if err != nil {
			if err == repository.ErrUserAlreadyExists {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "username already exists",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "failed to create user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		})
	}
}

func Login(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		user, err := userRepo.GetUserByUsername(c.Request.Context(), req.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid username or password",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid username or password",
			})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Username)
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
	}
}

func Me(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"user_id":  userID,
			"username": username,
		},
	})
}
