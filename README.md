# Golang Clean Architecture Boilerplate

A modern, production-ready backend boilerplate built with **Go (Golang)**, **Gin Web Framework**, **GORM ORM**, **PostgreSQL**, and **JWT Authentication**, strictly adhering to **Clean Architecture** principles.

---

## 🌟 Key Features

- **Clean Architecture Layout**: Strict separation of concerns (*Domain*, *Repository*, *Service/Usecase*, *Delivery/HTTP*).
- **Gin Web Framework**: High-performance HTTP routing with custom middleware.
- **GORM + PostgreSQL**: Feature-rich ORM for type-safe database operations.
- **Database Migrations**: Version-controlled database migrations powered by `golang-migrate` (`up`, `down`, `version`, `force`, `create`).
- **JWT Authentication & RBAC**: Access tokens, refresh tokens, and Role-Based Access Control (Admin / User).
- **Password Security**: Strong password hashing using `bcrypt`.
- **Structured Logging**: Structured, leveled logging using Go's built-in `log/slog`.
- **Standardized API Responses**: Consistent JSON response envelopes (`success`, `message`, `data`, `meta`, `errors`).
- **Graceful Shutdown**: OS signal interception (`SIGINT`/`SIGTERM`) to finish active requests before stopping.
- **Container Ready**: Production-grade multi-stage `Dockerfile` and `docker-compose.yml`.
- **Developer Utilities**: Comprehensive `Makefile` for build, test, migration, and development tasks.

---

## 📁 Project Structure

```text
.
├── cmd/
│   ├── api/
│   │   └── main.go                 # Application entry point & dependency wiring
│   └── migrate/
│       └── main.go                 # CLI runner for database migrations
├── config/
│   └── config.go                   # Environment configuration loader
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/            # HTTP Controllers (Auth, User, Health)
│   │       ├── middleware/         # Auth JWT, CORS, Logger, Panic Recovery
│   │       └── router/             # Gin route definitions
│   ├── domain/                     # Entities, DTOs, and Repository/Service Interfaces
│   ├── repository/
│   │   └── postgres/               # GORM database implementation
│   └── service/                    # Business Logic Layer
├── pkg/
│   ├── database/                   # Database connection pool & migration runner
│   ├── logger/                     # Structured Logger wrapper (slog)
│   ├── response/                   # Standardized JSON response helpers
│   └── security/                   # Bcrypt hashing & JWT token manager
├── migrations/                     # SQL migration files (000001_..., etc.)
├── .env.example                    # Environment variables template
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # PostgreSQL + App container stack
├── Makefile                        # Command shortcuts
└── README.md
```

---

## 🚀 Getting Started

### 1. Running with Docker (Recommended)

Make sure Docker & Docker Compose are installed, then run:

```bash
# Start the application and PostgreSQL database
make docker-up

# Stop containers (preserves database volume)
make docker-down
```

The API will be available at `http://localhost:8080`.

### 2. Running Locally

1. Duplicate `.env.example` to create your `.env` file:
   ```bash
   cp .env.example .env
   ```
2. Adjust your PostgreSQL database credentials in `.env`.
3. Apply database migrations:
   ```bash
   make migrate-up
   ```
4. Start the server:
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

---

## 📌 API Endpoints

| Method | Endpoint | Auth | Role | Description |
| :--- | :--- | :---: | :---: | :--- |
| `GET` | `/health` | ❌ | All | Service health check |
| `POST` | `/api/v1/auth/register` | ❌ | All | Register a new user account |
| `POST` | `/api/v1/auth/login` | ❌ | All | Authenticate & obtain JWT tokens |
| `POST` | `/api/v1/auth/refresh` | ❌ | All | Refresh Access Token using Refresh Token |
| `GET` | `/api/v1/auth/me` | 🔒 | All | Get profile of the authenticated user |
| `GET` | `/api/v1/users` | 🔒 | All | List users with pagination |
| `GET` | `/api/v1/users/:id` | 🔒 | All | Get user details by ID |
| `PUT` | `/api/v1/users/:id` | 🔒 | All | Update user details |
| `DELETE` | `/api/v1/users/:id` | 🔒 | Admin | Delete a user (Admin only) |

---

## 🗄️ Database Migrations

This boilerplate uses `golang-migrate` to ensure all schema changes are versioned, reproducible, and tracked.

| Command | Description |
| :--- | :--- |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback the last applied migration step (-1) |
| `make migrate-version` | Show current migration version & state (*clean* or *dirty*) |
| `make migrate-force v=1` | Force set database version (clears *dirty state* after a failed migration) |
| `make migrate-create name=<name>` | Create a new pair of sequential `.up.sql` and `.down.sql` migration files |

#### Creating a New Migration:
```bash
make migrate-create name=create_products_table
# Generates:
# migrations/000002_create_products_table.up.sql
# migrations/000002_create_products_table.down.sql
```

---

## 🧪 Running Tests

Execute unit tests and view coverage reports:

```bash
make test
```

---

## ⚙️ Environment Variables (`.env`)

| Variable | Default | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Environment mode (`development` / `production`) |
| `APP_PORT` | `8080` | HTTP Server port |
| `APP_NAME` | `boilerplate-golang` | Application identifier name |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `boilerplate_db` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode (`disable`, `require`, etc.) |
| `DB_AUTO_MIGRATE` | `false` | Automatically run migrations upon application startup |
| `JWT_SECRET` | `supersecretjwtkey...` | Secret key for signing JWT tokens |
| `JWT_EXPIRATION_HOURS` | `24` | Access Token expiration lifetime in hours |
| `JWT_REFRESH_EXPIRATION_DAYS` | `7` | Refresh Token expiration lifetime in days |
| `CORS_ALLOWED_ORIGINS` | `*` | Allowed CORS origins (comma-separated or `*`) |

---

## 📄 License

This project is licensed under the MIT License. You are free to use and adapt it for personal or commercial projects.
