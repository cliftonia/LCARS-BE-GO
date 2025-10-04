# Square Enix Backend API

A RESTful API backend built with Go for the Square Enix iOS application.

## Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: Standard HTTP methods and status codes
- **In-Memory Storage**: Mock data repositories for quick development
- **CORS Support**: Configurable cross-origin resource sharing
- **Middleware**: Logging, recovery, and error handling
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Health Check**: Built-in endpoint for monitoring

## Project Structure

```
subspace-backend/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/
│   ├── config/          # Configuration management
│   │   └── config.go
│   ├── domain/          # Domain models and interfaces
│   │   ├── user.go
│   │   └── message.go
│   ├── http/            # HTTP layer
│   │   ├── handlers/    # HTTP handlers
│   │   │   ├── user_handler.go
│   │   │   ├── message_handler.go
│   │   │   └── helpers.go
│   │   ├── middleware/  # HTTP middleware
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   └── router.go    # Route configuration
│   └── repository/      # Data access layer
│       └── memory/      # In-memory implementations
│           ├── user_repository.go
│           └── message_repository.go
├── .env.example         # Environment variables template
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Docker Compose configuration
├── Makefile            # Build automation
└── README.md           # This file
```

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

## Getting Started

### 1. Clone the repository

```bash
cd subspace-backend
```

### 2. Set up environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Install dependencies

```bash
make deps
```

### 4. Run the application

**Option A: Local development**
```bash
make run
```

**Option B: Docker**
```bash
make docker-up
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### Users
- `GET /api/v1/users` - List all users (supports pagination)
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Messages
- `GET /api/v1/messages/{id}` - Get message by ID
- `POST /api/v1/messages` - Create new message
- `DELETE /api/v1/messages/{id}` - Delete message
- `PATCH /api/v1/messages/{id}/read` - Mark message as read
- `GET /api/v1/users/{userId}/messages` - Get user's messages
- `GET /api/v1/users/{userId}/messages/unread-count` - Get unread count

## Example Requests

### Get all users
```bash
curl http://localhost:8080/api/v1/users
```

### Get user by ID
```bash
curl http://localhost:8080/api/v1/users/user-1
```

### Create a new user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New User",
    "email": "newuser@example.com"
  }'
```

### Get user messages
```bash
curl http://localhost:8080/api/v1/users/user-1/messages
```

### Get unread message count
```bash
curl http://localhost:8080/api/v1/users/user-1/messages/unread-count
```

## Development

### Run tests
```bash
make test
```

### Run with coverage
```bash
make test-coverage
```

### Format code
```bash
make fmt
```

### Run linter
```bash
make vet
```

### Build binary
```bash
make build
```

## Docker Commands

### Build Docker image
```bash
make docker-build
```

### Start containers
```bash
make docker-up
```

### Stop containers
```bash
make docker-down
```

### View logs
```bash
make docker-logs
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options:

- `PORT` - Server port (default: 8080)
- `HOST` - Server host (default: localhost)
- `ENVIRONMENT` - Environment mode (development/production)
- `CORS_ALLOWED_ORIGINS` - Comma-separated list of allowed origins

## Architecture

### Domain Layer
Contains business entities and repository interfaces. This layer has no dependencies on other layers.

### Repository Layer
Implements data access using the repository pattern. Currently uses in-memory storage for rapid development.

### HTTP Layer
Handles HTTP requests/responses, routing, and middleware. Clean separation between transport and business logic.

### Configuration
Centralized configuration management with validation and environment-specific defaults.

## Future Enhancements

- [ ] PostgreSQL database integration
- [ ] Authentication & authorization (JWT)
- [ ] Rate limiting
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Comprehensive test coverage
- [ ] Logging improvements (structured logging)
- [ ] Metrics and monitoring
- [ ] GraphQL support
- [ ] WebSocket support for real-time features

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
