# Auth API Documentation

## Overview
Authentication endpoints for user login, registration, logout, and token refresh.

## Endpoints

### POST /auth/login
Login with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "user",
    "verified": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "access_token": "jwt_token",
  "refresh_token": "refresh_token",
  "expires_at": "2024-01-01T01:00:00Z"
}
```

### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** Same as login response (201 Created)

### POST /auth/logout
Logout user and invalidate tokens.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body (optional):**
```json
{
  "refresh_token": "refresh_token"
}
```

**Response:**
```json
{
  "message": "Logout successful"
}
```

### POST /auth/refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh_token"
}
```

**Response:**
```json
{
  "access_token": "new_jwt_token",
  "refresh_token": "new_refresh_token",
  "expires_at": "2024-01-01T01:00:00Z"
}
```