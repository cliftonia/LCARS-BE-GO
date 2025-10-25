# Subspace Backend - Quick Start Guide

Get the backend running in under 5 minutes!

## Prerequisites Check

```bash
# Check if you have the required tools
go version          # Should be 1.24.0+
docker --version    # For containerized setup
psql --version      # For local PostgreSQL (optional)
```

## 🚀 Fastest Way: Docker Desktop (Recommended)

### 1. Start Docker Desktop

- **macOS**: Open Docker Desktop from Applications
- **CLI**: `open -a Docker` (wait ~20 seconds for it to start)
- Verify it's running: `docker ps`

### 2. Quick Start

```bash
# Create environment file
cp .env.example .env

# Start everything with one command
docker compose up -d

# Or if you have docker-compose:
docker-compose up -d

# View logs
docker compose logs -f
```

**That's it!** The API is now running at `http://localhost:8080`

### 3. Verify It's Working

```bash
# Health check
curl http://localhost:8080/health

# Register a test user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Stop Services

```bash
docker compose down          # Stop services
docker compose down -v       # Stop and remove volumes (fresh start)
```

---

## 🛠️ Alternative: Local Development Setup

If you prefer to run Go directly (better for development/debugging):

### 1. Setup Environment

```bash
cp .env.example .env
```

### 2. Start PostgreSQL (Choose One)

**Option A: Docker PostgreSQL Only**
```bash
docker run -d \
  --name subspace-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=subspace \
  -p 5432:5432 \
  -v "$(pwd)/scripts/init.sql:/docker-entrypoint-initdb.d/init.sql" \
  postgres:16-alpine
```

**Option B: Homebrew PostgreSQL**
```bash
brew install postgresql@16
brew services start postgresql@16
createdb subspace
psql subspace < scripts/init.sql
```

### 3. Run the Backend

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go

# Or build and run
make build
./bin/server
```

---

## 📦 Using Make Commands

The project includes a Makefile for common tasks:

```bash
make help              # Show all available commands
make deps              # Download Go dependencies
make build             # Build the binary
make run               # Run the application
make test              # Run tests
make docker-up         # Start with Docker Compose
make docker-down       # Stop Docker services
make docker-logs       # View container logs
```

---

## 🐛 Troubleshooting

### Docker daemon not running
```bash
# macOS
open -a Docker
# Wait 20 seconds, then try: docker ps
```

### Port 8080 already in use
```bash
lsof -i :8080                    # Find the process
kill -9 <PID>                    # Kill it
# Or change PORT in .env file
```

### Port 5432 already in use (PostgreSQL)
```bash
# Check what's using it
lsof -i :5432

# If it's an old PostgreSQL:
brew services stop postgresql@16
# Or:
docker stop subspace-db
```

### Database connection failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs subspace-db

# Test connection manually
psql -h localhost -U postgres -d subspace
# Password: postgres
```

### Permission denied errors
```bash
# Make sure scripts are executable
chmod +x scripts/*.sh

# Or run with explicit bash
bash scripts/setup.sh
```

---

## 🧪 Test Credentials (Development)

The `init.sql` creates these test accounts:

| Email | Password | User ID |
|-------|----------|---------|
| admin@subspace.dev | admin123 | 00000000-0000-0000-0000-000000000001 |
| test@subspace.dev | admin123 | 00000000-0000-0000-0000-000000000002 |

---

## 📱 Next Steps

1. **Test the API**: See [API_QUICK_REFERENCE.md](./API_QUICK_REFERENCE.md)
2. **Mobile Integration**: See [MOBILE_INTEGRATION.md](./MOBILE_INTEGRATION.md)
3. **Run Tests**: `go test -v ./...`
4. **View Coverage**: `make test-coverage && open coverage.html`

---

## 🔥 One-Line Setup (for the impatient)

```bash
cp .env.example .env && docker compose up -d && sleep 5 && curl http://localhost:8080/health
```

If Docker is not running, start it first:
```bash
open -a Docker && sleep 20 && cp .env.example .env && docker compose up -d && sleep 5 && curl http://localhost:8080/health
```

---

## 📊 API Endpoints Quick Reference

- **Health**: `GET /health`
- **Register**: `POST /api/v1/auth/register`
- **Login**: `POST /api/v1/auth/login`
- **Get User**: `GET /api/v1/auth/me` (requires auth)
- **List Users**: `GET /api/v1/users` (requires auth)
- **Messages**: `GET /api/v1/users/{userId}/messages` (requires auth)

Full documentation: [README.md](./README.md)
