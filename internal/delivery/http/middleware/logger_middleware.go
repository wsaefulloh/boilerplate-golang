package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lakejorge/boilerplate-golang/pkg/logger"
)

const RequestIDHeader = "X-Request-ID"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header(RequestIDHeader, requestID)

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			"request_id", requestID,
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency.String(),
		)
	}
}
