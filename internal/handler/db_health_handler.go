package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBHealthCheck(mysqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mysqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "database unavailable",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"database": "mysql",
				"status":   "ok",
			},
		})
	}
}
