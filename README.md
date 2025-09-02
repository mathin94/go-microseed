```markdown
# microseed — Go microservice skeleton

Small starter for building HTTP microservices in Go with embedded DB migrations and seeding.

## Prerequisites
- Go 1.21\+ recommended
- A running database compatible with your `goose` migrations
- `.env` file for local config

## Quick start
```bash
cp .env.example .env
go mod tidy
go run ./cmd/app serve
# healthz:  http://localhost:8080/healthz
# readyz:   http://localhost:8080/readyz
```

## CLI
```bash
# Run API
go run ./cmd/app serve

# DB migrations (embedded via goose)
go run ./cmd/app migrate up
go run ./cmd/app migrate down
go run ./cmd/app migrate reset

# Seed data
go run ./cmd/app seed
```

## Make usage
If you use the provided `Makefile`, common targets are:
```bash
# Show available targets
make help

# Install deps and prepare local env
make setup

# Run API (alias of serve)
make run
make serve

# DB migrations
make migrate-up
make migrate-down
make migrate-reset

# Seed data
make seed

# Lint, format, tidy, test, build
make lint
make fmt
make tidy
make test
make build
```

## Configuration
Provide local config via `.env`:
- `APP_PORT` \= default `8080`
- `DB_DSN` \= database connection string
- any other service-specific keys

```bash
# Example
APP_PORT=8080
DB_DSN=postgres://user:pass@localhost:5432/app?sslmode=disable
```

## Health endpoints
- `GET /healthz` \= liveness
- `GET /readyz` \= readiness
```