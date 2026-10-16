package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wsaefulloh/boilerplate-golang/pkg/response"
	"github.com/wsaefulloh/boilerplate-golang/pkg/security"
)

const (
	AuthorizationHeader = "Authorization"
	BearerSchema        = "Bearer "
	CtxUserIDKey        = "user_id"
	CtxUserEmailKey     = "user_email"
	CtxUserRoleKey      = "user_role"
)

func AuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerSchema) {
			response.Unauthorized(c, "Invalid authorization format, must be Bearer token")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, BearerSchema)
		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired access token")
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUserEmailKey, claims.Email)
		c.Set(CtxUserRoleKey, claims.Role)
		c.Next()
	}
}

func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(CtxUserRoleKey)
		if !exists {
			response.Unauthorized(c, "User role not found in context")
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			response.Forbidden(c, "Access denied: invalid role")
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == userRole {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Access denied: insufficient permissions")
		c.Abort()
	}
}
