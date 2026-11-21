# API examples and schema

Example requests against a local server, and the PostgreSQL schema. The full contract is in [`openapi.yaml`](openapi.yaml), also served at `GET /openapi.yaml`.

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
