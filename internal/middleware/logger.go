package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		start := time.Now()
		c.Next()

		latency := time.Since(start)
		log.Printf(
			"request_id=%s method=%s path=%s status=%d latency=%s client_ip=%s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			formatLatency(latency),
			c.ClientIP(),
		)
	}
}

func newRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func formatLatency(d time.Duration) string {
	return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1e6)
}
