# Subspace Mobile Integration Guide

Complete guide for integrating the Subspace Backend API with iOS and Android applications.

## Table of Contents

- [Quick Start](#quick-start)
- [Backend Setup](#backend-setup)
- [Authentication Flow](#authentication-flow)
- [API Endpoints](#api-endpoints)
- [iOS Integration](#ios-integration)
- [Android Integration](#android-integration)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Prerequisites

- Backend running at `http://localhost:8080` (or your deployment URL)
- Mobile device/simulator on the same network (for localhost testing)
- Understanding of REST APIs and JSON

### Base Configuration

```
Base URL: http://localhost:8080
API Version: v1
Content-Type: application/json
Authentication: Bearer Token (JWT)
```

## Backend Setup

### 1. Start the Backend

**Using Docker Compose (Recommended):**
```bash
cd subspace-backend
docker compose up -d
```

**Using Local Go:**
```bash
cd subspace-backend
go run cmd/server/main.go
```

### 2. Verify Backend is Running

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-10T00:00:00Z"
}
```

### 3. Access from Mobile Device

**Option A: Use ngrok (for localhost testing)**
```bash
# Install ngrok: https://ngrok.com/download
ngrok http 8080
```

This gives you a public URL like: `https://abc123.ngrok.io`

**Option B: Use your computer's local IP**
```bash
# Find your local IP
ifconfig | grep "inet " | grep -v 127.0.0.1

# Update CORS in .env
CORS_ALLOWED_ORIGINS=http://192.168.1.100:8080,http://localhost:8080
```

Then use `http://192.168.1.100:8080` as your base URL.

## Authentication Flow

### Step 1: User Registration

**Endpoint:** `POST /api/v1/auth/register`

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "createdAt": "2025-10-10T10:00:00Z",
    "updatedAt": "2025-10-10T10:00:00Z"
  }
}
```

**Validation Rules:**
- Name: Required, max 255 characters
- Email: Valid email format, max 255 characters, unique
- Password: Minimum 8 characters

### Step 2: User Login

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "createdAt": "2025-10-10T10:00:00Z",
    "updatedAt": "2025-10-10T10:00:00Z"
  }
}
```

### Step 3: Store the JWT Token

**Store securely:**
- iOS: Keychain Services
- Android: EncryptedSharedPreferences or Keystore

**Token Properties:**
- Expiration: 24 hours (configurable)
- Algorithm: HS256
- Claims: `user_id`, `email`, `exp`, `iat`

### Step 4: Use Token for Authenticated Requests

**Add to all protected endpoint requests:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Step 5: Get Current User

**Endpoint:** `GET /api/v1/auth/me`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "createdAt": "2025-10-10T10:00:00Z",
  "updatedAt": "2025-10-10T10:00:00Z"
}
```

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Health Check
```
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-10T10:00:00Z"
}
```

#### Register
```
POST /api/v1/auth/register
Content-Type: application/json

Body:
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

Body:
{
  "email": "string",
  "password": "string"
}
```

### Protected Endpoints (Authentication Required)

**All protected endpoints require:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Get Current User
```
GET /api/v1/auth/me
```

#### List All Users
```
GET /api/v1/users?limit=20&offset=0
```

**Query Parameters:**
- `limit` (optional): Items per page (default: 20, max: 100, min: 1)
- `offset` (optional): Number of items to skip (default: 0)

**Response:**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2025-10-10T10:00:00Z",
      "updatedAt": "2025-10-10T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

#### Get User by ID
```
GET /api/v1/users/{userId}
```

#### Create User
```
POST /api/v1/users
Content-Type: application/json

Body:
{
  "name": "string",
  "email": "string"
}
```

#### Update User
```
PUT /api/v1/users/{userId}
Content-Type: application/json

Body:
{
  "name": "string",
  "email": "string"
}
```

#### Delete User
```
DELETE /api/v1/users/{userId}
```

#### Get User Messages
```
GET /api/v1/users/{userId}/messages?limit=20&offset=0
```

**Response:**
```json
{
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "content": "Welcome to Subspace!",
      "kind": "info",
      "isRead": false,
      "createdAt": "2025-10-10T10:00:00Z",
      "updatedAt": "2025-10-10T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

#### Get Unread Message Count
```
GET /api/v1/users/{userId}/messages/unread-count
```

**Response:**
```json
{
  "count": 5
}
```

#### Get Message by ID
```
GET /api/v1/messages/{messageId}
```

#### Create Message
```
POST /api/v1/messages
Content-Type: application/json

Body:
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Message content here",
  "kind": "info"
}
```

**Message Kinds:**
- `info` - Informational message
- `warning` - Warning message
- `error` - Error message
- `success` - Success message

**Validation:**
- Content: Max 5000 characters

#### Mark Message as Read
```
PATCH /api/v1/messages/{messageId}/read
```

#### Delete Message
```
DELETE /api/v1/messages/{messageId}
```

## iOS Integration

### 1. Network Layer Setup

**Create an API Client:**

```swift
import Foundation

final class SubspaceAPIClient {
    static let shared = SubspaceAPIClient()

    private let baseURL = "http://localhost:8080/api/v1"
    private let session: URLSession

    private init() {
        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 30
        config.timeoutIntervalForResource = 300
        self.session = URLSession(configuration: config)
    }

    // MARK: - Authentication

    func register(name: String, email: String, password: String) async throws -> AuthResponse {
        let endpoint = "/auth/register"
        let body = RegisterRequest(name: name, email: email, password: password)
        return try await post(endpoint: endpoint, body: body)
    }

    func login(email: String, password: String) async throws -> AuthResponse {
        let endpoint = "/auth/login"
        let body = LoginRequest(email: email, password: password)
        return try await post(endpoint: endpoint, body: body)
    }

    func getCurrentUser(token: String) async throws -> User {
        let endpoint = "/auth/me"
        return try await get(endpoint: endpoint, token: token)
    }

    // MARK: - Users

    func getUsers(limit: Int = 20, offset: Int = 0, token: String) async throws -> PaginatedResponse<User> {
        let endpoint = "/users?limit=\(limit)&offset=\(offset)"
        return try await get(endpoint: endpoint, token: token)
    }

    func getUser(id: String, token: String) async throws -> User {
        let endpoint = "/users/\(id)"
        return try await get(endpoint: endpoint, token: token)
    }

    // MARK: - Messages

    func getUserMessages(userId: String, limit: Int = 20, offset: Int = 0, token: String) async throws -> PaginatedResponse<Message> {
        let endpoint = "/users/\(userId)/messages?limit=\(limit)&offset=\(offset)"
        return try await get(endpoint: endpoint, token: token)
    }

    func getUnreadCount(userId: String, token: String) async throws -> UnreadCountResponse {
        let endpoint = "/users/\(userId)/messages/unread-count"
        return try await get(endpoint: endpoint, token: token)
    }

    func markAsRead(messageId: String, token: String) async throws {
        let endpoint = "/messages/\(messageId)/read"
        try await patch(endpoint: endpoint, token: token)
    }

    func createMessage(userId: String, content: String, kind: MessageKind, token: String) async throws -> Message {
        let endpoint = "/messages"
        let body = CreateMessageRequest(userId: userId, content: content, kind: kind)
        return try await post(endpoint: endpoint, body: body, token: token)
    }

    // MARK: - Private Helpers

    private func get<T: Decodable>(endpoint: String, token: String? = nil) async throws -> T {
        var request = URLRequest(url: URL(string: baseURL + endpoint)!)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            if let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data) {
                throw APIError.serverError(errorResponse.error)
            }
            throw APIError.statusCode(httpResponse.statusCode)
        }

        return try JSONDecoder().decode(T.self, from: data)
    }

    private func post<T: Encodable, U: Decodable>(endpoint: String, body: T, token: String? = nil) async throws -> U {
        var request = URLRequest(url: URL(string: baseURL + endpoint)!)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        request.httpBody = try JSONEncoder().encode(body)

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            if let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data) {
                throw APIError.serverError(errorResponse.error)
            }
            throw APIError.statusCode(httpResponse.statusCode)
        }

        return try JSONDecoder().decode(U.self, from: data)
    }

    private func patch(endpoint: String, token: String) async throws {
        var request = URLRequest(url: URL(string: baseURL + endpoint)!)
        request.httpMethod = "PATCH"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            if let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data) {
                throw APIError.serverError(errorResponse.error)
            }
            throw APIError.statusCode(httpResponse.statusCode)
        }
    }
}

// MARK: - Models

struct User: Codable, Identifiable {
    let id: String
    let name: String
    let email: String
    let createdAt: Date
    let updatedAt: Date
}

struct Message: Codable, Identifiable {
    let id: String
    let userId: String
    let content: String
    let kind: MessageKind
    let isRead: Bool
    let createdAt: Date
    let updatedAt: Date
}

enum MessageKind: String, Codable {
    case info
    case warning
    case error
    case success
}

struct AuthResponse: Codable {
    let token: String
    let user: User
}

struct PaginatedResponse<T: Codable>: Codable {
    let data: [T]
    let total: Int
    let limit: Int
    let offset: Int
}

struct UnreadCountResponse: Codable {
    let count: Int
}

struct RegisterRequest: Codable {
    let name: String
    let email: String
    let password: String
}

struct LoginRequest: Codable {
    let email: String
    let password: String
}

struct CreateMessageRequest: Codable {
    let userId: String
    let content: String
    let kind: MessageKind
}

struct ErrorResponse: Codable {
    let error: String
}

enum APIError: LocalizedError {
    case invalidResponse
    case statusCode(Int)
    case serverError(String)

    var errorDescription: String? {
        switch self {
        case .invalidResponse:
            return "Invalid response from server"
        case .statusCode(let code):
            return "Server returned status code: \(code)"
        case .serverError(let message):
            return message
        }
    }
}
```

### 2. Token Storage (iOS)

**Create a Keychain Helper:**

```swift
import Foundation
import Security

final class KeychainHelper {
    static let shared = KeychainHelper()
    private let service = "com.subspace.app"

    private init() {}

    func save(_ token: String, for key: String) {
        let data = Data(token.utf8)

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]

        SecItemDelete(query as CFDictionary)
        SecItemAdd(query as CFDictionary, nil)
    }

    func get(for key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true
        ]

        var result: AnyObject?
        SecItemCopyMatching(query as CFDictionary, &result)

        guard let data = result as? Data else { return nil }
        return String(data: data, encoding: .utf8)
    }

    func delete(for key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key
        ]

        SecItemDelete(query as CFDictionary)
    }
}

// Usage
KeychainHelper.shared.save(token, for: "jwt_token")
let token = KeychainHelper.shared.get(for: "jwt_token")
```

### 3. SwiftUI Example

```swift
import SwiftUI

@MainActor
final class AuthViewModel: ObservableObject {
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var errorMessage: String?
    @Published var isLoading = false

    private let apiClient = SubspaceAPIClient.shared
    private let keychain = KeychainHelper.shared

    init() {
        checkAuthStatus()
    }

    func checkAuthStatus() {
        if let token = keychain.get(for: "jwt_token") {
            Task {
                do {
                    currentUser = try await apiClient.getCurrentUser(token: token)
                    isAuthenticated = true
                } catch {
                    keychain.delete(for: "jwt_token")
                    isAuthenticated = false
                }
            }
        }
    }

    func login(email: String, password: String) async {
        isLoading = true
        errorMessage = nil

        do {
            let response = try await apiClient.login(email: email, password: password)
            keychain.save(response.token, for: "jwt_token")
            currentUser = response.user
            isAuthenticated = true
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    func register(name: String, email: String, password: String) async {
        isLoading = true
        errorMessage = nil

        do {
            let response = try await apiClient.register(name: name, email: email, password: password)
            keychain.save(response.token, for: "jwt_token")
            currentUser = response.user
            isAuthenticated = true
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    func logout() {
        keychain.delete(for: "jwt_token")
        currentUser = nil
        isAuthenticated = false
    }
}

struct LoginView: View {
    @StateObject private var viewModel = AuthViewModel()
    @State private var email = ""
    @State private var password = ""

    var body: some View {
        VStack(spacing: 20) {
            TextField("Email", text: $email)
                .textInputAutocapitalization(.never)
                .keyboardType(.emailAddress)
                .textFieldStyle(.roundedBorder)

            SecureField("Password", text: $password)
                .textFieldStyle(.roundedBorder)

            if let error = viewModel.errorMessage {
                Text(error)
                    .foregroundColor(.red)
                    .font(.caption)
            }

            Button("Login") {
                Task {
                    await viewModel.login(email: email, password: password)
                }
            }
            .disabled(viewModel.isLoading || email.isEmpty || password.isEmpty)
        }
        .padding()
    }
}
```

## Android Integration

### 1. Add Dependencies (build.gradle.kts)

```kotlin
dependencies {
    // Retrofit for networking
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    implementation("com.squareup.okhttp3:logging-interceptor:4.11.0")

    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")

    // EncryptedSharedPreferences
    implementation("androidx.security:security-crypto:1.1.0-alpha06")
}
```

### 2. Create API Service

```kotlin
import retrofit2.Response
import retrofit2.http.*

interface SubspaceApiService {
    // Auth
    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<AuthResponse>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>

    @GET("auth/me")
    suspend fun getCurrentUser(@Header("Authorization") token: String): Response<User>

    // Users
    @GET("users")
    suspend fun getUsers(
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0,
        @Header("Authorization") token: String
    ): Response<PaginatedResponse<User>>

    @GET("users/{id}")
    suspend fun getUser(
        @Path("id") id: String,
        @Header("Authorization") token: String
    ): Response<User>

    // Messages
    @GET("users/{userId}/messages")
    suspend fun getUserMessages(
        @Path("userId") userId: String,
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0,
        @Header("Authorization") token: String
    ): Response<PaginatedResponse<Message>>

    @GET("users/{userId}/messages/unread-count")
    suspend fun getUnreadCount(
        @Path("userId") userId: String,
        @Header("Authorization") token: String
    ): Response<UnreadCountResponse>

    @POST("messages")
    suspend fun createMessage(
        @Body request: CreateMessageRequest,
        @Header("Authorization") token: String
    ): Response<Message>

    @PATCH("messages/{id}/read")
    suspend fun markAsRead(
        @Path("id") messageId: String,
        @Header("Authorization") token: String
    ): Response<Unit>
}

// Models
data class User(
    val id: String,
    val name: String,
    val email: String,
    val createdAt: String,
    val updatedAt: String
)

data class Message(
    val id: String,
    val userId: String,
    val content: String,
    val kind: MessageKind,
    val isRead: Boolean,
    val createdAt: String,
    val updatedAt: String
)

enum class MessageKind {
    info, warning, error, success
}

data class AuthResponse(
    val token: String,
    val user: User
)

data class PaginatedResponse<T>(
    val data: List<T>,
    val total: Int,
    val limit: Int,
    val offset: Int
)

data class UnreadCountResponse(
    val count: Int
)

data class RegisterRequest(
    val name: String,
    val email: String,
    val password: String
)

data class LoginRequest(
    val email: String,
    val password: String
)

data class CreateMessageRequest(
    val userId: String,
    val content: String,
    val kind: MessageKind
)
```

### 3. Create API Client

```kotlin
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

object SubspaceApiClient {
    private const val BASE_URL = "http://10.0.2.2:8080/api/v1/" // Android emulator localhost

    private val loggingInterceptor = HttpLoggingInterceptor().apply {
        level = HttpLoggingInterceptor.Level.BODY
    }

    private val okHttpClient = OkHttpClient.Builder()
        .addInterceptor(loggingInterceptor)
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    val api: SubspaceApiService = Retrofit.Builder()
        .baseUrl(BASE_URL)
        .client(okHttpClient)
        .addConverterFactory(GsonConverterFactory.create())
        .build()
        .create(SubspaceApiService::class.java)
}
```

### 4. Token Storage (Android)

```kotlin
import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey

class SecureTokenStorage(context: Context) {
    private val masterKey = MasterKey.Builder(context)
        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
        .build()

    private val sharedPreferences = EncryptedSharedPreferences.create(
        context,
        "subspace_secure_prefs",
        masterKey,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    fun saveToken(token: String) {
        sharedPreferences.edit().putString("jwt_token", token).apply()
    }

    fun getToken(): String? {
        return sharedPreferences.getString("jwt_token", null)
    }

    fun clearToken() {
        sharedPreferences.edit().remove("jwt_token").apply()
    }
}
```

### 5. Repository Pattern

```kotlin
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class AuthRepository(
    private val api: SubspaceApiService,
    private val tokenStorage: SecureTokenStorage
) {
    suspend fun register(name: String, email: String, password: String): Result<AuthResponse> {
        return withContext(Dispatchers.IO) {
            try {
                val response = api.register(RegisterRequest(name, email, password))
                if (response.isSuccessful && response.body() != null) {
                    val authResponse = response.body()!!
                    tokenStorage.saveToken(authResponse.token)
                    Result.success(authResponse)
                } else {
                    Result.failure(Exception(response.errorBody()?.string() ?: "Unknown error"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
    }

    suspend fun login(email: String, password: String): Result<AuthResponse> {
        return withContext(Dispatchers.IO) {
            try {
                val response = api.login(LoginRequest(email, password))
                if (response.isSuccessful && response.body() != null) {
                    val authResponse = response.body()!!
                    tokenStorage.saveToken(authResponse.token)
                    Result.success(authResponse)
                } else {
                    Result.failure(Exception(response.errorBody()?.string() ?: "Unknown error"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
    }

    suspend fun getCurrentUser(): Result<User> {
        return withContext(Dispatchers.IO) {
            try {
                val token = tokenStorage.getToken() ?: return@withContext Result.failure(Exception("No token"))
                val response = api.getCurrentUser("Bearer $token")
                if (response.isSuccessful && response.body() != null) {
                    Result.success(response.body()!!)
                } else {
                    Result.failure(Exception(response.errorBody()?.string() ?: "Unknown error"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
    }

    fun logout() {
        tokenStorage.clearToken()
    }

    fun isLoggedIn(): Boolean {
        return tokenStorage.getToken() != null
    }
}
```

### 6. ViewModel Example (Jetpack Compose)

```kotlin
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class AuthViewModel(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<AuthUiState>(AuthUiState.Idle)
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    private val _isAuthenticated = MutableStateFlow(false)
    val isAuthenticated: StateFlow<Boolean> = _isAuthenticated.asStateFlow()

    init {
        checkAuthStatus()
    }

    private fun checkAuthStatus() {
        _isAuthenticated.value = authRepository.isLoggedIn()
        if (_isAuthenticated.value) {
            viewModelScope.launch {
                authRepository.getCurrentUser()
                    .onSuccess { user ->
                        _uiState.value = AuthUiState.Success(user)
                    }
                    .onFailure {
                        authRepository.logout()
                        _isAuthenticated.value = false
                    }
            }
        }
    }

    fun login(email: String, password: String) {
        viewModelScope.launch {
            _uiState.value = AuthUiState.Loading
            authRepository.login(email, password)
                .onSuccess { authResponse ->
                    _isAuthenticated.value = true
                    _uiState.value = AuthUiState.Success(authResponse.user)
                }
                .onFailure { error ->
                    _uiState.value = AuthUiState.Error(error.message ?: "Unknown error")
                }
        }
    }

    fun register(name: String, email: String, password: String) {
        viewModelScope.launch {
            _uiState.value = AuthUiState.Loading
            authRepository.register(name, email, password)
                .onSuccess { authResponse ->
                    _isAuthenticated.value = true
                    _uiState.value = AuthUiState.Success(authResponse.user)
                }
                .onFailure { error ->
                    _uiState.value = AuthUiState.Error(error.message ?: "Unknown error")
                }
        }
    }

    fun logout() {
        authRepository.logout()
        _isAuthenticated.value = false
        _uiState.value = AuthUiState.Idle
    }
}

sealed class AuthUiState {
    object Idle : AuthUiState()
    object Loading : AuthUiState()
    data class Success(val user: User) : AuthUiState()
    data class Error(val message: String) : AuthUiState()
}
```

## Error Handling

### HTTP Status Codes

| Status Code | Meaning | Common Causes |
|------------|---------|---------------|
| 200 | OK | Successful GET/PATCH request |
| 201 | Created | Successful POST (registration, message creation) |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid input, validation error |
| 401 | Unauthorized | Missing/invalid/expired token |
| 404 | Not Found | Resource doesn't exist |
| 429 | Too Many Requests | Rate limit exceeded (100 req/min) |
| 500 | Internal Server Error | Server error |

### Error Response Format

```json
{
  "error": "Human-readable error message"
}
```

### Common Error Scenarios

**1. Invalid Email:**
```json
{
  "error": "invalid email format"
}
```

**2. Password Too Short:**
```json
{
  "error": "password must be at least 8 characters"
}
```

**3. Email Already Exists:**
```json
{
  "error": "email already exists"
}
```

**4. Invalid Credentials:**
```json
{
  "error": "invalid credentials"
}
```

**5. Token Expired:**
```json
{
  "error": "token has expired"
}
```

**6. Rate Limit Exceeded:**
```
HTTP 429 Too Many Requests
```

### Handling Token Expiration

**iOS Example:**
```swift
func handleAPIError(_ error: APIError) {
    switch error {
    case .statusCode(401):
        // Token expired, logout user
        KeychainHelper.shared.delete(for: "jwt_token")
        // Navigate to login screen
    case .serverError(let message) where message.contains("expired"):
        // Token expired
        KeychainHelper.shared.delete(for: "jwt_token")
    default:
        // Handle other errors
        break
    }
}
```

**Android Example:**
```kotlin
private fun handleError(error: Throwable) {
    when {
        error.message?.contains("401") == true -> {
            // Token expired, logout
            authRepository.logout()
            // Navigate to login
        }
        error.message?.contains("expired") == true -> {
            authRepository.logout()
        }
        else -> {
            // Handle other errors
        }
    }
}
```

## Best Practices

### 1. Secure Token Storage
- **iOS**: Use Keychain Services (never UserDefaults)
- **Android**: Use EncryptedSharedPreferences or Keystore

### 2. Network Security
- **Always use HTTPS in production**
- Add SSL pinning for extra security
- Implement certificate validation

### 3. Rate Limiting
- Implement exponential backoff for retry logic
- Respect 429 responses
- Cache data locally when possible

### 4. Error Handling
- Always handle 401 (redirect to login)
- Show user-friendly error messages
- Log errors for debugging (not in production)

### 5. Token Refresh
- Current implementation: 24-hour token expiration
- Re-authenticate when token expires
- Consider implementing refresh tokens for better UX

### 6. Offline Support
- Cache user data locally
- Queue messages when offline
- Sync when connection restored

### 7. Request Timeouts
- Set appropriate timeouts (30 seconds recommended)
- Show loading indicators
- Allow users to cancel long requests

### 8. Pagination
- Load data in chunks (default 20 items)
- Implement infinite scroll or pagination UI
- Track offset for "load more" functionality

## Testing

### 1. Test User Accounts

The database is pre-seeded with test accounts:

```
Email: admin@subspace.dev
Password: admin123

Email: test@subspace.dev
Password: admin123
```

### 2. Testing on Simulator/Emulator

**iOS Simulator:**
```
Base URL: http://localhost:8080
```

**Android Emulator:**
```
Base URL: http://10.0.2.2:8080
```

### 3. Testing on Physical Device

**Option A - Use ngrok:**
```bash
ngrok http 8080
# Use the provided HTTPS URL in your app
```

**Option B - Local network:**
```bash
# Find your computer's IP
ifconfig | grep "inet " | grep -v 127.0.0.1

# Update CORS in backend .env
CORS_ALLOWED_ORIGINS=http://192.168.1.100:8080

# Use http://192.168.1.100:8080 in your app
```

### 4. API Testing Tools

**Test with curl:**
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"password123"}'

# Get current user
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Test with Postman:**
1. Import collection from API endpoints above
2. Set base URL to `http://localhost:8080/api/v1`
3. Add Authorization header after login

## Troubleshooting

### Connection Refused

**Problem:** Can't connect to backend from mobile device

**Solutions:**
1. Verify backend is running: `curl http://localhost:8080/health`
2. Check firewall settings
3. Use correct IP for device type (simulator vs emulator vs physical)
4. Update CORS settings in backend `.env`

### 401 Unauthorized

**Problem:** All requests return 401

**Solutions:**
1. Verify token is being sent: `Authorization: Bearer TOKEN`
2. Check token hasn't expired (24 hours)
3. Ensure token was saved correctly after login
4. Re-login to get new token

### 429 Too Many Requests

**Problem:** Rate limit exceeded

**Solutions:**
1. Wait 1 minute before retrying
2. Implement request throttling in app
3. Cache data to reduce requests
4. Increase rate limit in backend (production only)

### CORS Errors

**Problem:** Requests blocked by CORS policy

**Solutions:**
1. Add your origin to `CORS_ALLOWED_ORIGINS` in backend `.env`
2. Restart backend after changing CORS settings
3. Use proper URL format: `http://192.168.1.100:8080`

### SSL/TLS Errors (Production)

**Problem:** Certificate validation fails

**Solutions:**
1. Ensure using HTTPS in production
2. Verify SSL certificate is valid
3. Update trust store if using self-signed certs
4. Check date/time on device

### Data Not Appearing

**Problem:** Empty responses or missing data

**Solutions:**
1. Check database has data: Access `/health` endpoint
2. Verify user ID in requests
3. Check pagination parameters
4. Look at server logs for errors

## Additional Resources

- **Backend Repository**: [GitHub Link]
- **API Documentation**: `http://localhost:8080/api/v1` (when Swagger is added)
- **Backend Logs**: `docker compose logs -f api`
- **Database Access**: `docker exec -it subspace-db psql -U postgres -d subspace`

## Support

For issues or questions:
1. Check this integration guide
2. Review backend README.md
3. Check backend logs for errors
4. Open an issue on GitHub

---

**Last Updated:** 2025-10-10
**Backend Version:** v1.0.0
**Minimum iOS:** 16.0
**Minimum Android:** API 26 (Android 8.0)
