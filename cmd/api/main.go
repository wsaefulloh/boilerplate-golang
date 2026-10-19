package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wsaefulloh/boilerplate-golang/config"
	"github.com/wsaefulloh/boilerplate-golang/internal/delivery/http/handler"
	"github.com/wsaefulloh/boilerplate-golang/internal/delivery/http/router"
	"github.com/wsaefulloh/boilerplate-golang/internal/repository/postgres"
	"github.com/wsaefulloh/boilerplate-golang/internal/service"
	"github.com/wsaefulloh/boilerplate-golang/pkg/database"
	"github.com/wsaefulloh/boilerplate-golang/pkg/logger"
	"github.com/wsaefulloh/boilerplate-golang/pkg/security"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// 2. Initialize Logger
	logger.InitLogger(cfg.App.Env)
	logger.Info("Starting application...",
		"app_name", cfg.App.Name,
		"env", cfg.App.Env,
		"port", cfg.App.Port,
	)

	// 3. Connect to Database
	isDev := cfg.App.Env == "development"
	db, err := database.NewPostgresDB(&cfg.Database, isDev)
	if err != nil {
		logger.Warn("Failed to connect to PostgreSQL database", "error", err.Error())
		logger.Warn("Application will continue running, but DB-dependent endpoints will fail until DB is available.")
	} else {
		// Run SQL migrations if enabled in configuration (DB_AUTO_MIGRATE=true)
		if cfg.Database.AutoMigrate {
			logger.Info("DB_AUTO_MIGRATE is enabled. Running database migrations...")
			if err := database.MigrateUp(&cfg.Database, "migrations"); err != nil {
				logger.Error("Failed to run database migrations", "error", err.Error())
			}
		}
	}

	// 4. Initialize Security Services
	jwtService := security.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.ExpirationHours,
		cfg.JWT.RefreshExpirationDays,
	)

	// 5. Initialize Repositories (handling nil db if database failed to connect initially)
	var userRepo = postgres.NewUserRepository(db)

	// 6. Initialize Services
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)

	// 7. Initialize Handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 8. Setup Router
	r := router.SetupRouter(&router.RouterDependencies{
		Config:        cfg,
		HealthHandler: healthHandler,
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		JWTService:    jwtService,
	})

	// 9. Configure HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 10. Start Server in Goroutine
	go func() {
		logger.Info("Server is listening and serving HTTP", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed to listen", "error", err.Error())
			os.Exit(1)
		}
	}()

	// 11. Graceful Shutdown Signal Handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err.Error())
	} else {
		logger.Info("Server exited cleanly")
	}
}
