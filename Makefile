.PHONY: run build test tidy docker-up docker-down docker-clean clean help migrate-up migrate-down migrate-version migrate-force migrate-create shell db-shell

APP_NAME = api
BUILD_DIR = bin

help: ## Display available commands
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

run: ## Run the application locally
	@go run cmd/api/main.go

build: ## Build the application binary
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME) & $(BUILD_DIR)/migrate"

test: ## Run unit tests
	@go test -v -cover ./...

tidy: ## Tidy dependencies
	@go mod tidy
	@go mod verify

migrate-up: ## Run all pending database migrations
	@go run cmd/migrate/main.go up

migrate-down: ## Rollback 1 migration step
	@go run cmd/migrate/main.go down

migrate-version: ## Check current database migration version and status
	@go run cmd/migrate/main.go version

migrate-force: ## Force set migration version (e.g. make migrate-force v=1)
	@go run cmd/migrate/main.go force $(v)

migrate-create: ## Create new migration files (e.g. make migrate-create name=add_orders)
	@go run cmd/migrate/main.go create $(name)

docker-up: ## Start containers with docker-compose
	@docker-compose up -d --build

docker-down: ## Stop docker-compose containers (preserves database volume)
	@docker-compose down

docker-clean: ## Stop containers and remove volumes (WARNING: wipes database)
	@docker-compose down -v

shell: ## Open interactive shell inside app container
	@docker exec -it boilerplate_app sh

db-shell: ## Open psql interactive shell inside postgres container
	@docker exec -it boilerplate_postgres psql -U postgres -d boilerplate_db

clean: ## Clean binary files
	@rm -rf $(BUILD_DIR)
