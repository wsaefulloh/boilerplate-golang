package database

import (
	"fmt"
	"time"

	"github.com/wsaefulloh/boilerplate-golang/config"
	"github.com/wsaefulloh/boilerplate-golang/internal/domain"
	"github.com/wsaefulloh/boilerplate-golang/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DBConfig, isDev bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
	)

	gormLogLevel := gormlogger.Error
	if isDev {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Connection pool configuration
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database auto-migrations...")
	err := db.AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	logger.Info("Database auto-migration completed successfully")
	return nil
}
