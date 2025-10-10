# Subspace API Quick Reference

Quick reference for mobile developers integrating with the Subspace Backend API.

## Base Configuration

```
Base URL: http://localhost:8080/api/v1
Content-Type: application/json
Authorization: Bearer {token}
```

## Authentication Flow

```
1. Register/Login → Receive JWT token
2. Store token securely (Keychain/EncryptedSharedPreferences)
3. Add token to all protected requests: Authorization: Bearer {token}
4. Token expires after 24 hours → Re-authenticate
```

## Quick Start URLs

| Environment | iOS Simulator | Android Emulator | Physical Device |
|------------|---------------|------------------|-----------------|
| Localhost | `http://localhost:8080` | `http://10.0.2.2:8080` | `http://{YOUR_IP}:8080` |
| Ngrok | `https://{id}.ngrok.io` | `https://{id}.ngrok.io` | `https://{id}.ngrok.io` |

## Endpoints Cheat Sheet

### 🔓 Public (No Auth Required)

```http
# Health Check
GET /health

# Register
POST /api/v1/auth/register
{ "name": "string", "email": "string", "password": "string" }

# Login
POST /api/v1/auth/login
{ "email": "string", "password": "string" }
```

### 🔒 Protected (Auth Required)

```http
# Get Current User
GET /api/v1/auth/me
Authorization: Bearer {token}

# List Users (paginated)
GET /api/v1/users?limit=20&offset=0
Authorization: Bearer {token}

# Get User by ID
GET /api/v1/users/{id}
Authorization: Bearer {token}

# Get User Messages (paginated)
GET /api/v1/users/{userId}/messages?limit=20&offset=0
Authorization: Bearer {token}

# Get Unread Count
GET /api/v1/users/{userId}/messages/unread-count
Authorization: Bearer {token}

# Create Message
POST /api/v1/messages
Authorization: Bearer {token}
{ "userId": "string", "content": "string", "kind": "info" }

# Mark Message as Read
PATCH /api/v1/messages/{id}/read
Authorization: Bearer {token}
```

## Response Formats

### Auth Response
```json
{
  "token": "eyJhbGci...",
  "user": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "createdAt": "timestamp",
    "updatedAt": "timestamp"
  }
}
```

### Paginated Response
```json
{
  "data": [...],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

## Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| 200 | OK | Success |
| 201 | Created | Resource created |
| 400 | Bad Request | Check input validation |
| 401 | Unauthorized | Re-authenticate |
| 404 | Not Found | Resource doesn't exist |
| 429 | Too Many Requests | Wait 1 minute |
| 500 | Server Error | Check backend logs |

## Validation Rules

```
Email: Valid format, max 255 chars, unique
Password: Min 8 chars
Name: Max 255 chars
Message Content: Max 5000 chars
Message Kind: info | warning | error | success
Pagination Limit: 1-100, default 20
Pagination Offset: ≥ 0, default 0
```

## Rate Limiting

```
100 requests per minute per IP
Burst: 10 requests
Response when exceeded: HTTP 429
```

## iOS Code Snippet

```swift
// Network Request
var request = URLRequest(url: URL(string: baseURL + "/auth/login")!)
request.httpMethod = "POST"
request.setValue("application/json", forHTTPHeaderField: "Content-Type")
request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
request.httpBody = try JSONEncoder().encode(body)

let (data, response) = try await URLSession.shared.data(for: request)

// Store Token
KeychainHelper.shared.save(token, for: "jwt_token")
let token = KeychainHelper.shared.get(for: "jwt_token")
```

## Android Code Snippet

```kotlin
// Retrofit Service
interface SubspaceApiService {
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>

    @GET("users/{id}")
    suspend fun getUser(
        @Path("id") id: String,
        @Header("Authorization") token: String
    ): Response<User>
}

// Store Token
val tokenStorage = SecureTokenStorage(context)
tokenStorage.saveToken(token)
val token = tokenStorage.getToken()
```

## Message Kinds

```kotlin
enum MessageKind {
    info     // Blue - Informational
    warning  // Yellow - Warning
    error    // Red - Error
    success  // Green - Success
}
```

## Test Accounts

```
Email: admin@subspace.dev
Password: admin123

Email: test@subspace.dev
Password: admin123
```

## Common Headers

```http
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Troubleshooting Quick Fixes

| Problem | Solution |
|---------|----------|
| Connection refused | Check backend running, use correct IP |
| 401 Unauthorized | Check token format: `Bearer {token}` |
| CORS error | Add your origin to backend CORS config |
| 429 Rate limit | Wait 60 seconds, implement throttling |
| SSL error | Use HTTPS in production |

## cURL Examples

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get current user
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Get messages
curl "http://localhost:8080/api/v1/users/USER_ID/messages?limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Backend Commands

```bash
# Start backend
docker compose up -d

# View logs
docker compose logs -f api

# Stop backend
docker compose down

# Check health
curl http://localhost:8080/health
```

## Next Steps

1. ✅ Start backend: `docker compose up -d`
2. ✅ Test health: `curl http://localhost:8080/health`
3. ✅ Register user via API
4. ✅ Store JWT token securely
5. ✅ Make authenticated requests
6. ✅ Handle errors (especially 401)
7. ✅ Implement token refresh logic

---

**For detailed integration guide, see:** [MOBILE_INTEGRATION.md](./MOBILE_INTEGRATION.md)
