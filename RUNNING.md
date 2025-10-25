# 🎉 Your Subspace Backend is Running!

## Quick Status Check

```bash
# Check if containers are running
docker ps

# View logs
docker compose logs -f

# Health check
curl http://localhost:8080/health
```

## 🔥 Test the API

### 1. Register a New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Your Name",
    "email": "you@example.com",
    "password": "yourpassword"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "ea1eebef-3957-4675-b6f5-c9556eebd19c",
    "name": "Your Name",
    "email": "you@example.com",
    ...
  }
}
```

### 2. Get Current User (Authenticated)
```bash
# Save your token from registration
TOKEN="your-token-here"

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

### 3. List All Users
```bash
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Create a Message
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "your-user-id",
    "content": "Hello from Subspace!",
    "kind": "info"
  }'
```

### 5. Get User Messages
```bash
curl http://localhost:8080/api/v1/users/YOUR_USER_ID/messages \
  -H "Authorization: Bearer $TOKEN"
```

## 📊 Available Endpoints

### Public Endpoints
- `GET /health` - Health check
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

### Protected Endpoints (Require JWT Token)
- `GET /api/v1/auth/me` - Get current user
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/messages/{id}` - Get message by ID
- `POST /api/v1/messages` - Create message
- `DELETE /api/v1/messages/{id}` - Delete message
- `PATCH /api/v1/messages/{id}/read` - Mark message as read
- `GET /api/v1/users/{userId}/messages` - Get user's messages
- `GET /api/v1/users/{userId}/messages/unread-count` - Get unread count

## 🗄️ Database Access

```bash
# Connect to PostgreSQL
docker exec -it subspace-db psql -U postgres -d subspace

# Check users
docker exec subspace-db psql -U postgres -d subspace -c "SELECT * FROM users;"

# Check messages
docker exec subspace-db psql -U postgres -d subspace -c "SELECT * FROM messages;"
```

## 🎮 Control Commands

```bash
# View logs (all services)
docker compose logs -f

# View API logs only
docker compose logs -f api

# View database logs only
docker compose logs -f postgres

# Restart services
docker compose restart

# Stop services
docker compose down

# Stop and remove all data
docker compose down -v

# Rebuild and restart
docker compose up -d --build
```

## 🧪 Run Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/http/handlers
go test -v ./internal/repository/postgres
```

## 🐛 Troubleshooting

### Services not running
```bash
docker compose ps
docker compose logs
```

### Port conflicts
```bash
# Check what's using port 8080
lsof -i :8080

# Or use a different port (edit .env)
PORT=8081
docker compose down && docker compose up -d
```

### Database issues
```bash
# Check database is healthy
docker exec subspace-db pg_isready -U postgres

# Reset database
docker compose down -v
docker compose up -d
```

### View full logs
```bash
# All logs since start
docker compose logs

# Follow logs in real-time
docker compose logs -f

# Last 100 lines
docker compose logs --tail=100
```

## 📱 Next Steps

1. **Mobile Integration**: Check [MOBILE_INTEGRATION.md](./MOBILE_INTEGRATION.md)
2. **API Reference**: Check [API_QUICK_REFERENCE.md](./API_QUICK_REFERENCE.md)
3. **Full Documentation**: Check [README.md](./README.md)

## 🔒 Note About Test Users

The database includes pre-seeded test users:
- `admin@subspace.dev` / `admin123`
- `test@subspace.dev` / `admin123`

**Note:** The password hash in `init.sql` may not match the bcrypt implementation. It's safer to register new users via the `/auth/register` endpoint.

## 🎯 Your Backend is Now Running at:
**http://localhost:8080**

Happy coding! 🚀
