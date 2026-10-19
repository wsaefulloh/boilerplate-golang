package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/wsaefulloh/boilerplate-golang/config"
	"github.com/wsaefulloh/boilerplate-golang/pkg/database"
	"github.com/wsaefulloh/boilerplate-golang/pkg/logger"
)

const defaultMigrationsPath = "migrations"

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  up              Apply all pending migrations")
	fmt.Println("  down            Rollback the last migration step (-1)")
	fmt.Println("  version         Display the current migration version and dirty status")
	fmt.Println("  force <version> Force set the migration version (useful for clearing dirty state)")
	fmt.Println("  create <name>   Create a new pair of up and down SQL migration files")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go down")
	fmt.Println("  go run cmd/migrate/main.go version")
	fmt.Println("  go run cmd/migrate/main.go force 1")
	fmt.Println("  go run cmd/migrate/main.go create add_products_table")
}

func main() {
	var migrationsDir string
	flag.StringVar(&migrationsDir, "dir", defaultMigrationsPath, "path to migrations folder")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	// Handle 'create' without connecting to DB
	if command == "create" {
		if len(args) < 2 {
			logger.Error("Migration name required. Example: go run cmd/migrate/main.go create create_products_table")
			os.Exit(1)
		}
		name := args[1]
		upFile, downFile, err := database.CreateMigration(migrationsDir, name)
		if err != nil {
			logger.Error("Failed to create migration", "error", err)
			os.Exit(1)
		}
		logger.Info("Created migration files", "dir", migrationsDir, "up", upFile, "down", downFile)
		return
	}

	// For database migration commands, load DB config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.InitLogger(cfg.App.Env)

	switch command {
	case "up":
		if err := database.MigrateUp(&cfg.Database, migrationsDir); err != nil {
			logger.Error("Migration up failed", "error", err)
			os.Exit(1)
		}
		logger.Info("Database migration completed successfully")

	case "down":
		if err := database.MigrateDown(&cfg.Database, migrationsDir); err != nil {
			logger.Error("Migration down failed", "error", err)
			os.Exit(1)
		}
		logger.Info("Database rollback completed successfully")

	case "version", "status":
		ver, dirty, err := database.MigrateVersion(&cfg.Database, migrationsDir)
		if err != nil {
			logger.Error("Failed to retrieve migration version", "error", err)
			os.Exit(1)
		}
		if ver == 0 {
			logger.Info("Current migration version: 0 (no migrations applied yet)")
		} else {
			if dirty {
				logger.Warn("Database state is DIRTY. Please fix the issue and run 'force <version>'", "version", ver)
			} else {
				logger.Info("Database state: CLEAN", "version", ver)
			}
		}

	case "force":
		if len(args) < 2 {
			logger.Error("Version number required. Example: go run cmd/migrate/main.go force 1")
			os.Exit(1)
		}
		ver, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Error("Invalid version number: must be an integer", "input", args[1])
			os.Exit(1)
		}
		if err := database.MigrateForce(&cfg.Database, migrationsDir, ver); err != nil {
			logger.Error("Failed to force version", "error", err)
			os.Exit(1)
		}
		logger.Info("Database version forced", "version", ver)

	default:
		logger.Error("Unknown command", "command", command)
		printUsage()
		os.Exit(1)
	}
}
