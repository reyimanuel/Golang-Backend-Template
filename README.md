# Go Backend Template

This repository is a simple REST API template built with Go, Gin, GORM, and PostgreSQL.
It is structured so feature code lives under `internal/api`, shared infrastructure lives under `internal/infrastructures`, and schema setup lives under `internal/migration`.

## Stack

- Go 1.25.6
- Gin for HTTP routing
- GORM for ORM and migrations
- PostgreSQL as the default database
- JWT using RSA key pairs
- godotenv for local environment loading

## Requirements

Before running this project, prepare the following:

- Go 1.25.6 or a compatible Go 1.25.x version
- PostgreSQL running locally or on a reachable host
- A `.env` file copied from `.env.example`
- RSA key files for JWT signing and validation in the `keys` folder

Minimum environment variables:

- `PORT`: HTTP server port
- `IS_PRODUCTION`: `true` or `false`
- `DB_USER`
- `DB_PASS`
- `DB_NAME`
- `DB_HOST`
- `DB_PORT`
- `DB_TIME_ZONE`
- `ALLOW_ORIGIN`
- `BASE_URL`
- `ACCESS_TOKEN_LIFE_TIME`
- `REFRESH_TOKEN_LIFE_TIME`
- `PRIVATE_KEY`
- `PUBLIC_KEY`
- `FRONTEND_URL`

## Getting Started

1. Install dependencies:

```bash
go mod tidy
```

2. Copy environment configuration:

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

3. Make sure the key paths in `.env` point to valid files in `keys/`.

4. Make sure PostgreSQL is available and the database in `.env` already exists.

5. Run migration and seed data:

```bash
go run ./cmd/server migrate
```

6. Start the API server:

```bash
go run ./cmd/server
```

## CLI Commands

The application entrypoint supports a small CLI from `cmd/server/main.go`.

- `go run ./cmd/server`: start the HTTP server
- `go run ./cmd/server migrate`: run schema migration and seeding
- `go run ./cmd/server migrate --force`: allow seeding outside localhost
- `go run ./cmd/server migrate-only`: run schema migration only
- `go run ./cmd/server reset`: drop tables, recreate schema, and reseed locally
- `go run ./cmd/server reset --force`: bypass localhost protection for reset

## Default Seeder Account

The migration seeder creates one default account if it does not already exist:

- Username: `admin`
- Email: `admin@example.com`
- Password: `admin123`

This seed is idempotent. Running migration again will not create duplicate accounts.

## Project Structure

### Root folders

#### `cmd/`

Application entrypoints live here.

Requirement:
Each runnable program should have its own subfolder under `cmd`. This template currently uses `cmd/server` as the HTTP and CLI entrypoint.

#### `internal/`

Private application code lives here. Code in this folder is intended only for this repository and should not be imported by external projects.

Requirement:
Place all business logic, handlers, configuration, middleware, and persistence code here rather than in the root.

#### `keys/`

Stores RSA private and public key files used for JWT signing and verification.

Requirement:
This folder must contain the key files referenced by `PRIVATE_KEY` and `PUBLIC_KEY` in `.env`. Do not commit real production keys.

#### `public/`

Stores publicly accessible generated files and static assets.

Requirement:
Put exported files, generated images, or reusable document assets here when they are meant to be served or accessed outside the application runtime.

#### `tmp/`

Temporary runtime or generated files can be stored here.

Requirement:
Use this folder only for ephemeral files that are safe to regenerate or delete.

## Internal Structure

### `internal/api/`

Contains HTTP-facing feature modules.

Requirement:
Each feature should usually get its own subfolder containing route registration, handler logic, request and response DTOs, service logic, and repository access.

Current contents:

- `router.go`: central route registration for the API
- `users/`: example feature module for user-related endpoints

### `internal/api/users/`

Example feature module showing how to organize a resource.

Requirement:
Keep transport, use-case, and data-access responsibilities separated by file.

Expected file roles:

- `routes.go`: define feature routes and attach middleware if needed
- `handler.go`: parse requests and write HTTP responses
- `dto.go`: request and response payload definitions
- `service.go`: business rules and orchestration
- `repository.go`: database queries and persistence access

### `internal/constants/`

Holds shared constant values used across the application.

Requirement:
Only put constants here when they are reused broadly or represent stable application-wide values.

### `internal/infrastructures/`

Contains technical building blocks shared by multiple features.

Requirement:
Put framework integration, configuration loading, database setup, authentication helpers, and common middleware here. Avoid feature-specific business logic in this layer.

#### `internal/infrastructures/config/`

Loads and validates environment-based application configuration.

Requirement:
Every config value required at startup should be defined and resolved here.

#### `internal/infrastructures/database/`

Database connection and database bootstrap code live here.

Requirement:
Keep connection initialization, pool setup, and database client wiring here instead of scattering it across handlers or services.

#### `internal/infrastructures/middleware/`

Shared Gin middleware such as CORS and authentication guards.

Requirement:
Only middleware with cross-cutting concerns should live here.

#### `internal/infrastructures/pkg/`

Shared internal helper packages.

Requirement:
Use this folder for reusable technical helpers that do not belong to a single feature.

Current helper groups:

- `errs/`: shared error helpers or error response structures
- `helpers/`: general utility helpers such as password hashing and validation
- `token/`: JWT generation, loading, and validation logic
- `utils/`: small common helper functions used across the codebase

### `internal/migration/`

Database models, schema migration, and seed logic live here.

Requirement:
Whenever a new table is added, update the model definitions and register the model in `Models`. Any default data creation should also be handled here.

Current file roles:

- `models.go`: GORM model declarations and migration registration list
- `migration.go`: schema migration, reset, and drop orchestration
- `seed.go`: default data seeding

### `internal/server/`

Application server bootstrap logic lives here.

Requirement:
Keep HTTP server initialization, middleware wiring, and route registration startup code here.

## Suggested Workflow For New Features

1. Add or update database models in `internal/migration/models.go`.
2. Add a feature folder under `internal/api/your-feature`.
3. Create `dto.go`, `routes.go`, `handler.go`, `service.go`, and `repository.go` as needed.
4. Register the new routes in `internal/api/router.go`.
5. Run `go run ./cmd/server migrate` if schema changes were made.

## Notes

- Seeding is blocked when `DB_HOST` is not `localhost` or `127.0.0.1`, unless `--force` is passed.
- `migrate-only` skips seeding completely.
- The template expects key files to exist before startup because config loading fails fast when key paths are missing.