# Backend Features Analysis - Subspace iOS App

## Executive Summary

Based on analysis of the Subspace iOS app, the backend is **mostly complete** but is missing several critical OAuth/Social authentication endpoints and WebSocket support for real-time features.

**Current Status:**
- ✅ Basic Auth (Login/Register) - **IMPLEMENTED**
- ✅ User Management - **IMPLEMENTED**
- ✅ Messages CRUD - **IMPLEMENTED**
- ❌ Apple Sign In - **MISSING**
- ❌ Google Sign In - **MISSING**
- ❌ Token Refresh - **MISSING**
- ❌ WebSocket Support - **MISSING**
- ❌ Logout Endpoint - **MISSING (optional)**
- ⚠️  User Avatar URLs - **PARTIALLY MISSING**

---

## 1. Current iOS App Features

### Authentication Features
The iOS app supports **4 authentication methods**:

1. **Email/Password Login** ✅
   - Endpoint: `POST /api/v1/auth/login`
   - Request: `{email, password}`
   - Response: `{user, token}`

2. **Email/Password Registration** ✅
   - Endpoint: `POST /api/v1/auth/register`
   - Request: `{name, email, password}`
   - Response: `{user, token}`

3. **Apple Sign In** ❌ **NOT IMPLEMENTED**
   - Expected: `POST /api/v1/auth/apple`
   - Request: `{userId, identityToken, authorizationCode, email, fullName}`
   - Response: `{user, token}`

4. **Google Sign In** ❌ **NOT IMPLEMENTED**
   - Expected: `POST /api/v1/auth/google`
   - Request: `{userId, idToken, accessToken, email, fullName}`
   - Response: `{user, token}`

### User Management ✅
- `GET /api/v1/auth/me` - Get current user
- `GET /api/v1/users` - List users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Message Features ✅
- `GET /api/v1/users/{userId}/messages` - Get user messages
- `GET /api/v1/users/{userId}/messages/unread-count` - Get unread count
- `GET /api/v1/messages/{id}` - Get message by ID
- `POST /api/v1/messages` - Create message
- `DELETE /api/v1/messages/{id}` - Delete message
- `PATCH /api/v1/messages/{id}/read` - Mark as read

### Token Management ❌
- **Expected:** `POST /api/v1/auth/refresh`
- **Request:** `{refreshToken}`
- **Response:** `{accessToken, refreshToken, expiresAt}`
- **Status:** NOT IMPLEMENTED

### Logout ❌
- **Expected:** `POST /api/v1/auth/logout`
- **Purpose:** Server-side token invalidation
- **Status:** NOT IMPLEMENTED (optional - client-side logout works)

### WebSocket Real-Time Updates ❌
The iOS app has a `WebSocketManager` that expects:
- **Endpoint:** `ws://localhost:8080/ws?userId={userId}`
- **Message Types:**
  - `new_message` - New message created
  - `message_read` - Message marked as read
  - `message_deleted` - Message deleted
- **Status:** NOT IMPLEMENTED

---

## 2. Missing Backend Features (Priority Order)

### 🔴 HIGH PRIORITY

#### 1. Apple Sign In Endpoint
**Why:** Core authentication method for iOS
**Endpoint:** `POST /api/v1/auth/apple`
**Request:**
```json
{
  "userId": "string",
  "identityToken": "string",
  "authorizationCode": "string",
  "email": "string?",
  "fullName": {
    "givenName": "string?",
    "familyName": "string?"
  }
}
```
**Response:**
```json
{
  "user": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "avatarURL": "string?",
    "createdAt": "ISO8601"
  },
  "token": "JWT"
}
```

**Implementation Notes:**
- Verify Apple identity token with Apple's servers
- Create user if doesn't exist (first-time sign in)
- Return existing user if already registered
- Generate JWT token

#### 2. Google Sign In Endpoint
**Why:** Alternative social auth method
**Endpoint:** `POST /api/v1/auth/google`
**Request:**
```json
{
  "userId": "string",
  "idToken": "string",
  "accessToken": "string",
  "email": "string",
  "fullName": "string?"
}
```
**Response:** Same as Apple Sign In

**Implementation Notes:**
- Verify Google ID token with Google's servers
- Create/fetch user
- Generate JWT token

#### 3. Token Refresh Endpoint
**Why:** Maintain user session without re-login
**Endpoint:** `POST /api/v1/auth/refresh`
**Request:**
```json
{
  "refreshToken": "string"
}
```
**Response:**
```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "expiresAt": "ISO8601"
}
```

**Current Backend Issue:**
The backend currently only returns a single `token` in auth responses. Need to:
1. Implement separate access + refresh tokens
2. Update auth response to include both tokens
3. Implement token refresh logic

### 🟡 MEDIUM PRIORITY

#### 4. WebSocket Support for Real-Time Updates
**Why:** Real-time message notifications without polling
**Endpoint:** WebSocket at `/ws`
**Query Params:** `userId={userId}`

**Events to Broadcast:**
```typescript
// Client receives these events
{
  "type": "new_message",
  "data": {
    "messageId": "uuid",
    "userId": "uuid",
    "content": "string",
    "kind": "info|warning|error|success",
    "createdAt": "ISO8601"
  }
}

{
  "type": "message_read",
  "data": {
    "messageId": "uuid",
    "userId": "uuid"
  }
}

{
  "type": "message_deleted",
  "data": {
    "messageId": "uuid"
  }
}
```

**Implementation:**
- Use `gorilla/websocket` package
- Authenticate WebSocket connections via JWT query param or header
- Maintain connection pool per user
- Broadcast events when messages are created/read/deleted

### 🟢 LOW PRIORITY

#### 5. User Avatar Support
**Status:** Database schema has no `avatar_url` field
**iOS App Expects:** `avatarURL` in User model

**Options:**
1. Add `avatar_url` column to users table
2. Implement avatar upload endpoint
3. Use Gravatar URLs based on email
4. Keep as optional feature (iOS handles missing avatars gracefully)

#### 6. Logout Endpoint (Optional)
**Endpoint:** `POST /api/v1/auth/logout`
**Why:** Server-side token invalidation, audit logging
**Status:** Client can logout locally by deleting tokens

---

## 3. Backend API Response Format Issues

### Current Backend Response
```json
{
  "token": "jwt-token",
  "user": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "createdAt": "ISO8601",
    "updatedAt": "ISO8601"
  }
}
```

### iOS App Expects
```json
{
  "user": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "avatarURL": "string?",  // ← Missing
    "createdAt": "ISO8601"
  },
  "token": "jwt-token"
}
```

**Differences:**
1. ❌ No `avatarURL` field in User response
2. ✅ Token format is correct (single JWT)
3. ⚠️  iOS expects separate access/refresh tokens for production
4. ❌ `updatedAt` not used by iOS app

---

## 4. Implementation Roadmap

### Phase 1: OAuth Authentication (1-2 days)
- [ ] Add Apple Sign In endpoint
  - [ ] Validate Apple identity token
  - [ ] Create/fetch user logic
  - [ ] Return JWT token
- [ ] Add Google Sign In endpoint
  - [ ] Validate Google ID token
  - [ ] Create/fetch user logic
  - [ ] Return JWT token
- [ ] Add unit tests for OAuth flows
- [ ] Update documentation

### Phase 2: Token Management (1 day)
- [ ] Implement dual-token system (access + refresh)
- [ ] Add `refresh_tokens` table
- [ ] Implement refresh endpoint
- [ ] Update auth responses to include both tokens
- [ ] Update middleware to handle token refresh
- [ ] Add token expiration validation

### Phase 3: WebSocket Real-Time (2-3 days)
- [ ] Add WebSocket handler with gorilla/websocket
- [ ] Implement connection authentication
- [ ] Create connection pool manager
- [ ] Broadcast events on message operations
- [ ] Add reconnection logic
- [ ] Add unit/integration tests

### Phase 4: Avatar Support (1 day)
- [ ] Add `avatar_url` column to users table
- [ ] Update User model
- [ ] Add avatar upload endpoint (S3/local storage)
- [ ] Add avatar validation
- [ ] Optional: Integrate Gravatar fallback

### Phase 5: Logout Enhancement (Optional, 0.5 days)
- [ ] Implement logout endpoint
- [ ] Add token blacklist (Redis recommended)
- [ ] Add audit logging

---

## 5. Recommended Next Steps

### Immediate (This Week)
1. **Add Apple Sign In** - Critical for iOS app store requirements
2. **Add Google Sign In** - Popular alternative auth method
3. **Implement Token Refresh** - Better UX, security

### Short Term (Next 2 Weeks)
4. **WebSocket Support** - Real-time features, better UX
5. **Avatar Upload** - Enhanced user profiles

### Long Term (Nice to Have)
6. **Logout Endpoint** - Security audit compliance
7. **Token Blacklist** - Advanced security
8. **Rate Limiting per User** - Currently only per-IP

---

## 6. Technology Recommendations

### OAuth Verification Libraries
**Apple Sign In:**
```go
// Use official Apple libraries or JWT verification
github.com/golang-jwt/jwt/v5
```

**Google Sign In:**
```go
google.golang.org/api/oauth2/v2
```

### WebSocket
```go
github.com/gorilla/websocket  // Already may be available
```

### Token Blacklist (Optional)
```go
github.com/go-redis/redis/v9  // For token blacklist
```

---

## 7. Database Schema Changes Needed

### Add Avatar Support
```sql
ALTER TABLE users
ADD COLUMN avatar_url VARCHAR(500);
```

### Add Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    INDEX idx_refresh_tokens_token (token),
    INDEX idx_refresh_tokens_user_id (user_id)
);
```

### Add OAuth Providers Tracking (Optional)
```sql
CREATE TABLE oauth_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('apple', 'google', 'email')),
    provider_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id),
    INDEX idx_oauth_providers_user_id (user_id)
);
```

---

## 8. Testing Considerations

### OAuth Testing
- Use Apple/Google test credentials
- Mock token validation in tests
- Test user creation vs. existing user flow

### WebSocket Testing
- Test connection authentication
- Test broadcast to multiple clients
- Test reconnection logic
- Test message delivery guarantees

### Token Refresh Testing
- Test expired token refresh
- Test invalid refresh token
- Test token rotation
- Test concurrent refresh requests

---

## Summary

**Total Missing Features:** 6 major items
**Estimated Implementation Time:** 5-8 days
**Critical Path:** OAuth → Token Refresh → WebSockets → Avatars

The backend is **70% feature-complete** for the iOS app. The main gaps are OAuth authentication and real-time features via WebSockets.
