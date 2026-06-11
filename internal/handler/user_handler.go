package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/auth"
	"github.com/hawk-roy/Night-Hawk/internal/model"
	"github.com/hawk-roy/Night-Hawk/internal/repository"
	"github.com/hawk-roy/Night-Hawk/internal/response"
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
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request")

			return
		}

		if req.Username == "" {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "username is required")
			return
		}

		if req.Password == "" {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "password is required")
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "failed to hash password")
			return
		}

		user, err := userRepo.CreateUser(c.Request.Context(), req.Username, string(passwordHash))
		if err != nil {
			if err == repository.ErrUserAlreadyExists {
				response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "username already exists")
				return
			}

			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "failed to create user")
			return
		}

		response.Success(c, gin.H{
			"id":       user.ID,
			"username": user.Username,
		})

	}
}

func Login(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request")
			return
		}

		if req.Username == "" {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "username is required")
			return
		}

		if req.Password == "" {
			response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "password is required")
			return
		}

		user, err := userRepo.GetUserByUsername(c.Request.Context(), req.Username)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "invalid username or password")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "invalid username or password")
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Username)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "failed to generate token")
			return
		}

		response.Success(c, gin.H{
			"token": token,
		})
	}
}

func Me(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
		return
	}

	username, ok := c.Get("username")
	if !ok {
		response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
		return
	}

	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
