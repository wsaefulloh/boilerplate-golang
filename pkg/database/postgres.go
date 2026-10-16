package database

import (
	"fmt"
	"log"
	"time"

	"github.com/lakejorge/boilerplate-golang/config"
	"github.com/lakejorge/boilerplate-golang/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DBConfig, isDev bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
	)

	gormLogLevel := logger.Error
	if isDev {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
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

	log.Println("[INFO] Successfully connected to PostgreSQL database")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("[INFO] Running database auto-migrations...")
	err := db.AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	log.Println("[INFO] Database auto-migration completed successfully")
	return nil
}
