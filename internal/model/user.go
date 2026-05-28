package model

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
}
