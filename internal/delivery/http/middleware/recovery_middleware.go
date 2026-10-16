package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/wsaefulloh/boilerplate-golang/pkg/logger"
	"github.com/wsaefulloh/boilerplate-golang/pkg/response"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logger.Error("Panic recovered",
					"error", fmt.Sprintf("%v", err),
					"stack", stack,
				)
				response.InternalServerError(c, "An unexpected internal server error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
