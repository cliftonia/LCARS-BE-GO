# LCARS-BE-GO

The Go REST API behind [Subspace](https://github.com/cliftonia/Subspace), an LCARS-themed iOS messaging app.

## Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **PostgreSQL Database**: Database with connection pooling
- **JWT Authentication**: Token-based authentication with bcrypt password hashing
- **Rate Limiting**: Per-IP rate limiting (100 req/min) to prevent API abuse
- **Request Tracing**: UUID request IDs attached to every log line
- **Structured Logging**: JSON logging with contextual information using Go's `log/slog`
- **Input Validation**: Validation of request bodies (email, password and content length)
- **CORS Support**: Configurable cross-origin resource sharing
- **Middleware Stack**: Authentication, logging, recovery, rate limiting, and request ID
- **Docker Support**: Multi-stage builds with Docker Compose for PostgreSQL
- **CI/CD Pipeline**: GitHub Actions with testing, linting, building, and security scanning
- **Health Checks**: Database connectivity and service health monitoring

## Mobile App Integration

**For iOS & Android developers:**
- **[Mobile Integration Guide](./MOBILE_INTEGRATION.md)** - Complete integration guide with code examples
- **[API Quick Reference](./API_QUICK_REFERENCE.md)** - Quick reference cheat sheet

These guides include:
- Complete authentication flow with JWT
- iOS (SwiftUI) code examples with async/await
- Android (Kotlin/Jetpack Compose) code examples with coroutines
- Secure token storage (Keychain/EncryptedSharedPreferences)
- Network layer setup (URLSession/Retrofit)
- Error handling patterns
- Testing on simulators, emulators, and physical devices

## Tech Stack

- **Language**: Go 1.24.0
- **Web Framework**: Gorilla Mux
- **Database**: PostgreSQL 16 with `lib/pq` driver
- **Authentication**: JWT (`golang-jwt/jwt/v5`) + Bcrypt
- **Rate Limiting**: `golang.org/x/time/rate`
- **Logging**: Go standard library `log/slog`
- **Testing**: Go standard testing package
- **CI/CD**: GitHub Actions
- **Containerization**: Docker + Docker Compose

## Project Structure

```
LCARS-BE-GO/
├── cmd/
│   └── server/               # Application entry point
│       └── main.go
├── internal/
│   ├── auth/                 # Authentication & password hashing
│   │   ├── jwt.go           # JWT token generation/validation
│   │   └── password.go      # Bcrypt password hashing
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── constants/            # Application constants
│   │   └── constants.go
│   ├── database/             # Database connection & pooling
│   │   └── postgres.go
│   ├── domain/               # Domain models and interfaces
│   │   ├── errors.go        # Typed domain errors
│   │   ├── user.go
│   │   └── message.go
│   ├── http/                 # HTTP layer
│   │   ├── handlers/         # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── message_handler.go
│   │   │   ├── validation.go
│   │   │   └── helpers.go
│   │   ├── middleware/       # HTTP middleware
│   │   │   ├── auth.go      # JWT authentication
│   │   │   ├── logger.go    # Request logging
│   │   │   ├── rate_limit.go
│   │   │   ├── recovery.go  # Panic recovery
│   │   │   └── request_id.go
│   │   └── router.go         # Route configuration
│   ├── logger/               # Logging setup
│   │   └── logger.go
│   └── repository/           # Data access layer
│       ├── memory/           # In-memory implementations (testing)
│       │   ├── user_repository.go
│       │   └── message_repository.go
│       └── postgres/         # PostgreSQL implementations
│           ├── user_repository.go
│           └── message_repository.go
├── scripts/
│   └── init.sql             # PostgreSQL database schema
├── .github/
│   └── workflows/
│       └── ci.yml           # CI/CD pipeline
├── .env.example             # Environment variables template
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # PostgreSQL + API services
├── Makefile                 # Build automation
└── README.md                # This file
```

## Prerequisites

- Go 1.24.0 or higher
- PostgreSQL 16 (or use Docker Compose)
- Docker & Docker Compose (optional, for containerized deployment)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/cliftonia/LCARS-BE-GO.git
cd LCARS-BE-GO
```

### 2. Set up environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```

**Key Environment Variables:**
```bash
# Server
PORT=8080
HOST=localhost
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=subspace
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# Security
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Rate Limiting
API_RATE_LIMIT=100  # requests per minute
```

### 3. Option A: Run with Docker Compose (Recommended)

```bash
# Start PostgreSQL and API server
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

The API will be available at `http://localhost:8080`

### 3. Option B: Run Locally

**Install dependencies:**
```bash
go mod download
```

**Start PostgreSQL:**
```bash
# Using Docker
docker run -d \
  --name subspace-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=subspace \
  -p 5432:5432 \
  postgres:16-alpine

# Initialize database
psql -h localhost -U postgres -d subspace -f scripts/init.sql
```

**Run the application:**
```bash
go run cmd/server/main.go
```

## API Endpoints

### Public Endpoints

#### Health Check
- `GET /health` - Check API and database health

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/me` - Get current user (requires authentication)

### Protected Endpoints (Require JWT Token)

#### Users
- `GET /api/v1/users` - List all users (supports pagination)
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

#### Messages
- `GET /api/v1/messages/{id}` - Get message by ID
- `POST /api/v1/messages` - Create new message
- `DELETE /api/v1/messages/{id}` - Delete message
- `PATCH /api/v1/messages/{id}/read` - Mark message as read
- `GET /api/v1/users/{userId}/messages` - Get user's messages
- `GET /api/v1/users/{userId}/messages/unread-count` - Get unread count

## Example Requests

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Test User",
    "email": "test@example.com",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Get current user (authenticated)
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get all users (authenticated)
```bash
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get user messages (authenticated)
```bash
curl http://localhost:8080/api/v1/users/user-id/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get unread message count (authenticated)
```bash
curl http://localhost:8080/api/v1/users/user-id/messages/unread-count \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Development

### Run tests
```bash
go test -v ./...
```

### Run tests with coverage
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### View coverage in browser
```bash
go tool cover -html=coverage.out
```

### Format code
```bash
go fmt ./...
```

### Run linter
```bash
go vet ./...
```

### Build binary
```bash
go build -o bin/server ./cmd/server
```

### Run with live reload (requires `air`)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## Testing

Test coverage by layer:

- **Repository Layer**: 99.2% coverage
- **Handler Layer**: 27% coverage
- **Authentication**: Registration, login, and JWT validation
- **Middleware**: Logger, recovery, rate limiting, auth

**Run specific test packages:**
```bash
# Test handlers
go test -v ./internal/http/handlers

# Test repositories
go test -v ./internal/repository/...

# Test with race detection
go test -v -race ./...
```

## CI/CD Pipeline

GitHub Actions workflow includes:

1. **Test Job**: Unit tests with race detection and coverage reporting
2. **Lint Job**: `golangci-lint` for code quality
3. **Build Job**: Compile binary and upload artifact
4. **Security Job**: `gosec` security scanner

**Pipeline runs on:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`

## Architecture

### Domain Layer
Contains business entities and repository interfaces. This layer has no dependencies on other layers and defines:
- Domain models (`User`, `Message`)
- Repository interfaces
- Typed domain errors

### Repository Layer
Implements data access using the repository pattern with:
- PostgreSQL implementations with connection pooling
- In-memory implementations for testing
- Context-based query timeouts (3 seconds)
- Transaction support ready

### HTTP Layer
Handles HTTP requests/responses, routing, and middleware:
- Clean separation between transport and business logic
- Middleware: auth, logging, recovery, rate limiting, request ID
- Input validation for all endpoints
- Proper HTTP status codes and error responses

### Authentication Layer
JWT-based authentication with:
- Token generation with configurable expiration
- Token validation and claims extraction
- Bcrypt password hashing (cost factor 12)
- Protected route middleware

### Configuration
Centralized configuration management with:
- Environment variable support
- Validation on startup
- Environment-specific defaults
- Docker Compose integration

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Messages Table
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK (LENGTH(content) > 0 AND LENGTH(content) <= 5000),
    kind VARCHAR(50) NOT NULL CHECK (kind IN ('info', 'warning', 'error', 'success')),
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `users.email` - Fast email lookups
- `messages.user_id` - Fast user message queries
- `messages.created_at` - Chronological ordering
- `messages.is_read` - Unread message filtering

## Connection Pooling

PostgreSQL connection pool configuration:
- **Max Open Connections**: 25
- **Max Idle Connections**: 5
- **Connection Max Lifetime**: 5 minutes
- **Connection Max Idle Time**: 5 minutes
- **Query Timeout**: 3 seconds per operation
- **Health Check**: 2-second timeout

## Rate Limiting

API rate limiting per IP address:
- **Limit**: 100 requests per minute
- **Burst**: 10 requests
- **Cleanup**: Auto-cleanup inactive limiters after 3 minutes
- **Response**: 429 Too Many Requests when exceeded

## Security Features

- **Password Hashing**: Bcrypt with cost factor 12
- **JWT Tokens**: HS256 algorithm, configurable expiration
- **Input Validation**: Email format, password length, content length
- **SQL Injection Protection**: Parameterized queries
- **XSS Protection**: Proper content type headers
- **CORS**: Configurable allowed origins
- **Panic Recovery**: Graceful error handling
- **Secrets Management**: Environment variables, never in code

## Monitoring & Observability

- **Structured Logging**: JSON format in production, text in development
- **Request IDs**: UUID per request, included in logs
- **Health Endpoints**: Database connectivity checks
- **HTTP Logging**: Method, path, status, duration, bytes, user agent

## Production Deployment

### Environment Variables (Production)

```bash
ENVIRONMENT=production
PORT=8080
HOST=0.0.0.0

# Use strong secrets in production
JWT_SECRET=your-super-secret-key-min-32-chars
DB_PASSWORD=strong-database-password

# Enable SSL for database
DB_SSL_MODE=require

# Restrict CORS
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

### Docker Build

```bash
# Build production image
docker build -t subspace-backend:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e ENVIRONMENT=production \
  -e JWT_SECRET=$JWT_SECRET \
  -e DB_HOST=$DB_HOST \
  subspace-backend:latest
```

### Health Check Endpoint

Use `/health` for:
- Load balancer health checks
- Kubernetes liveness/readiness probes
- Monitoring systems

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Performance

- **Database Queries**: 3-second timeout prevents hanging
- **Connection Pool**: Reuses connections efficiently
- **Memory**: Minimal allocations with proper cleanup
- **Graceful Shutdown**: 30-second timeout for in-flight requests

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs subspace-db

# Test connection
psql -h localhost -U postgres -d subspace
```

### Build Issues

```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod download
go mod tidy
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`go test -v ./...`)
4. Run linter (`go vet ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review test files for usage examples

## Roadmap

- [x] PostgreSQL database integration
- [x] Authentication & authorization (JWT)
- [x] Rate limiting
- [x] Repository-layer test coverage
- [x] Structured logging
- [x] CI/CD pipeline
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Metrics and monitoring (Prometheus)
- [ ] GraphQL support
- [ ] WebSocket support for real-time features
- [ ] Database migrations tool
- [ ] Admin dashboard
- [ ] Email notifications
