package database

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/wsaefulloh/boilerplate-golang/config"
)

// GetDatabaseURL constructs the postgres connection string for golang-migrate
func GetDatabaseURL(cfg *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}

// NewMigrator initializes a new golang-migrate instance
func NewMigrator(cfg *config.DBConfig, migrationsPath string) (*migrate.Migrate, error) {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))
	dbURL := GetDatabaseURL(cfg)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %w", err)
	}

	return m, nil
}

// MigrateUp applies all available up migrations
func MigrateUp(cfg *config.DBConfig, migrationsPath string) error {
	m, err := NewMigrator(cfg, migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	log.Printf("[INFO] Applying migrations from %s ...\n", migrationsPath)
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("[INFO] No new migrations to apply (database is up to date)")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	log.Println("[INFO] Successfully applied all migrations")
	return nil
}

// MigrateDown rolls back the last applied migration
func MigrateDown(cfg *config.DBConfig, migrationsPath string) error {
	m, err := NewMigrator(cfg, migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	log.Printf("[INFO] Rolling back 1 migration step from %s ...\n", migrationsPath)
	if err := m.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("[INFO] No migrations to rollback")
			return nil
		}
		return fmt.Errorf("migration down failed: %w", err)
	}

	log.Println("[INFO] Successfully rolled back 1 migration step")
	return nil
}

// MigrateVersion returns current migration version and dirty state
func MigrateVersion(cfg *config.DBConfig, migrationsPath string) (uint, bool, error) {
	m, err := NewMigrator(cfg, migrationsPath)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, err
	}

	return version, dirty, nil
}

// MigrateForce forces a specific migration version (cleans dirty state)
func MigrateForce(cfg *config.DBConfig, migrationsPath string, version int) error {
	m, err := NewMigrator(cfg, migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	log.Printf("[INFO] Forcing migration version to %d ...\n", version)
	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version %d: %w", version, err)
	}

	log.Printf("[INFO] Successfully forced migration version to %d\n", version)
	return nil
}

// CreateMigration creates a pair of .up.sql and .down.sql migration files
func CreateMigration(migrationsPath, name string) (upFile string, downFile string, err error) {
	if strings.TrimSpace(name) == "" {
		return "", "", errors.New("migration name cannot be empty")
	}

	// Sanitize name: lowercase and replace spaces with underscores
	sanitizedName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))

	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Find the next sequential version
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read migrations directory: %w", err)
	}

	maxVersion := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		parts := strings.Split(entry.Name(), "_")
		if len(parts) >= 2 {
			if v, err := strconv.Atoi(parts[0]); err == nil {
				if v > maxVersion {
					maxVersion = v
				}
			}
		}
	}

	nextVersion := maxVersion + 1
	versionPrefix := fmt.Sprintf("%06d", nextVersion)

	upFileName := fmt.Sprintf("%s_%s.up.sql", versionPrefix, sanitizedName)
	downFileName := fmt.Sprintf("%s_%s.down.sql", versionPrefix, sanitizedName)

	upPath := filepath.Join(migrationsPath, upFileName)
	downPath := filepath.Join(migrationsPath, downFileName)

	upContent := fmt.Sprintf("-- Migration: %s (up)\n", sanitizedName)
	downContent := fmt.Sprintf("-- Migration: %s (down)\n", sanitizedName)

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write up migration: %w", err)
	}

	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write down migration: %w", err)
	}

	return upFileName, downFileName, nil
}
