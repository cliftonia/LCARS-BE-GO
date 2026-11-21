# LCARS-BE-GO

[![CI](https://github.com/cliftonia/LCARS-BE-GO/actions/workflows/ci.yml/badge.svg)](https://github.com/cliftonia/LCARS-BE-GO/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.24-00ADD8)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791)

The Go REST API behind [Subspace](https://github.com/cliftonia/Subspace), an LCARS-themed iOS messaging app.
Accounts, sign-in (email, Apple, Google), refresh tokens and messages, backed by PostgreSQL.

## What's in it

- **Auth**: email and password with bcrypt, Sign in with Apple and Google token verification, short-lived JWT
  access tokens with rotating refresh tokens
- **Users and messages**: CRUD, per-user message lists, read state and unread counts
- **Operational basics**: structured `log/slog` logging with request IDs, per-IP rate limiting, panic recovery, CORS,
  a `/health` check that pings the database, and Prometheus metrics at `/metrics`
- **Contract**: an [OpenAPI spec](docs/openapi.yaml), also served at `GET /openapi.yaml`
- **Schema changes**: versioned SQL [migrations](migrations/), applied on start-up when `AUTO_MIGRATE=true`
- **CI**: tests with the race detector, `golangci-lint`, a build and a `gosec` scan on every push

## Architecture

```
cmd/server            wiring: config, database, optional migrations, router
internal/domain       entities, repository interfaces, typed errors  (no dependencies)
internal/repository   postgres/ for real, memory/ for tests          (implements domain interfaces)
internal/http         router, handlers, middleware                   (depends on domain only)
internal/auth         JWT, bcrypt, Apple and Google token checks
```

Handlers depend on the repository interfaces, so the HTTP tests run against the in-memory repositories with no
database.

## Endpoints

All under `/api/v1` except the health, metrics and spec routes.

| Method | Path | Auth |
|---|---|---|
| `POST` | `/auth/register`, `/auth/login`, `/auth/apple`, `/auth/google`, `/auth/refresh` | – |
| `GET` | `/auth/me` | JWT |
| `GET` `POST` | `/users` | JWT |
| `GET` `PUT` `DELETE` | `/users/{id}` | JWT |
| `GET` | `/users/{userId}/messages`, `/users/{userId}/messages/unread-count` | JWT |
| `POST` | `/messages` | JWT |
| `GET` `DELETE` | `/messages/{id}` | JWT |
| `PATCH` | `/messages/{id}/read` | JWT |
| `GET` | `/health`, `/metrics`, `/openapi.yaml` | – |

Example requests and the table definitions are in [`docs/api.md`](docs/api.md).

## Running it

```bash
cp .env.example .env      # set JWT_SECRET at least
docker compose up -d      # PostgreSQL and the API on localhost:8080
curl localhost:8080/health
```

Or against your own PostgreSQL: `go run ./cmd/server`. More in [`docs/running.md`](docs/running.md).

## Development

```bash
go test -race ./...       # unit and HTTP tests, no database needed
golangci-lint run ./...
make help                 # build, docker and migration targets
```

## Docs

- [Quick start](docs/quickstart.md)
- [Running and configuration](docs/running.md)
- [Mobile integration guide](docs/mobile-integration.md): iOS and Android client examples
- [API quick reference](docs/api-quick-reference.md)
- [API examples and schema](docs/api.md)

## Licence

[MIT](LICENSE)
