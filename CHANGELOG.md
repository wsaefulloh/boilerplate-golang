# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-10-17

### Added
- **Clean Architecture Implementation**: Layered structure separating Domain, Delivery (HTTP), Service (Business Logic), and Repository (Data Access).
- **Gin Web Framework**: High-performance HTTP router with custom middlewares (CORS, Request Logger, JWT Authentication, and Panic Recovery).
- **GORM with PostgreSQL**: Connection pooling and type-safe query interface.
- **Database Migrations (`golang-migrate`)**: Versioned schema migrations via CLI (`cmd/migrate/main.go`) and package helpers (`pkg/database/migrate.go`) supporting `up`, `down`, `version`, `force`, and `create`.
- **JWT Authentication & RBAC**: Access Token and Refresh Token lifecycle management with Role-Based Access Control (`admin` / `user`).
- **Password Security**: Password hashing and verification using `bcrypt`.
- **Structured Logging**: Production-grade structured logging using Go standard library `log/slog` (text handler for development, JSON handler for production).
- **Standardized API Responses**: Envelope structure (`success`, `message`, `data`, `meta`, `errors`) via `pkg/response`.
- **Graceful Shutdown**: OS signal interception (`SIGINT`/`SIGTERM`) ensuring non-disruptive server termination.
- **Docker Support**: Optimized multi-stage `Dockerfile` and `docker-compose.yml` for local containerization.
- **Developer Utilities**: `Makefile` with commands for building, running, testing, migrations, and container operations.
- **Initial Migration Files**: `000001_create_users_table.up.sql` and `000001_create_users_table.down.sql`.
- **Unit Tests**: Security tests covering password hashing and JWT issuance/validation.
