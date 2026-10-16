package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lakejorge/boilerplate-golang/config"
	"github.com/lakejorge/boilerplate-golang/internal/delivery/http/handler"
	"github.com/lakejorge/boilerplate-golang/internal/delivery/http/middleware"
	"github.com/lakejorge/boilerplate-golang/pkg/security"
)

type RouterDependencies struct {
	Config        *config.Config
	HealthHandler *handler.HealthHandler
	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandler
	JWTService    *security.JWTService
}

func SetupRouter(deps *RouterDependencies) *gin.Engine {
	if deps.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware(deps.Config.CORS.AllowedOrigins))

	// Health Check Route
	r.GET("/health", deps.HealthHandler.HealthCheck)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		// Public Auth Routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", deps.AuthHandler.Register)
			auth.POST("/login", deps.AuthHandler.Login)
			auth.POST("/refresh", deps.AuthHandler.RefreshToken)
		}

		// Authenticated Routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(deps.JWTService))
		{
			// Auth Profile
			protected.GET("/auth/me", deps.AuthHandler.GetMe)

			// User Management Routes
			users := protected.Group("/users")
			{
				users.GET("", deps.UserHandler.List)
				users.GET("/:id", deps.UserHandler.GetByID)
				users.PUT("/:id", deps.UserHandler.Update)
				
				// Admin-only route for deletion
				users.DELETE("/:id", middleware.RoleGuard("admin"), deps.UserHandler.Delete)
			}
		}
	}

	return r
}
