package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

func DBHealthCheck(mysqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mysqlDB.Ping(); err != nil {
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "database unavailable")
			return
		}

		response.Success(c, gin.H{
			"database": "mysql",
			"status":   "ok",
		})
	}
}
