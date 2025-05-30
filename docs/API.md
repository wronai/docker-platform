# 📡 API Documentation

This document provides comprehensive documentation for the Media Vault API, including endpoints, request/response formats, and authentication details.

## Table of Contents
- [Authentication](#-authentication)
- [Base URL](#-base-url)
- [Response Format](#-response-format)
- [Error Handling](#-error-handling)
- [Endpoints](#-endpoints)
  - [Authentication](#authentication-1)
  - [Users](#users)
  - [Media](#media)
  - [Collections](#collections)
  - [Sharing](#sharing)
  - [Search](#search)
- [Webhooks](#-webhooks)
- [Rate Limiting](#-rate-limiting)
- [Pagination](#-pagination)
- [WebSocket API](#-websocket-api)

## 🔑 Authentication

All API endpoints require authentication using JWT (JSON Web Tokens).

### Obtaining a Token

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "yourpassword"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### Using the Token

Include the token in the `Authorization` header:
```
Authorization: Bearer your.jwt.token.here
```

## 🌐 Base URL

All API endpoints are relative to the base URL:
```
https://api.your-mediavault-instance.com/api/v1
```

## 📦 Response Format

All successful API responses follow this format:

```json
{
  "success": true,
  "data": {
    /* response data */
  },
  "meta": {
    "timestamp": "2023-05-30T10:00:00Z",
    "version": "1.0.0"
  }
}
```

## ❌ Error Handling

Error responses include an error code and message:

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "validation_error",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Must be a valid email address"
      }
    ]
  },
  "meta": {
    "timestamp": "2023-05-30T10:00:00Z",
    "request_id": "req_1234567890"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `unauthorized` | 401 | Invalid or missing authentication |
| `forbidden` | 403 | Insufficient permissions |
| `not_found` | 404 | Resource not found |
| `validation_error` | 400 | Invalid input data |
| `rate_limited` | 429 | Too many requests |
| `internal_error` | 500 | Server error |

## 📡 Endpoints

### Authentication

#### Login
```http
POST /auth/login
```

#### Refresh Token
```http
POST /auth/refresh
```

#### Logout
```http
POST /auth/logout
```

### Users

#### Get Current User
```http
GET /users/me
```

#### Update Profile
```http
PATCH /users/me
```

### Media

#### Upload File
```http
POST /media/upload
Content-Type: multipart/form-data
```

#### Get Media Item
```http
GET /media/{id}
```

#### List Media
```http
GET /media
```

### Collections

#### Create Collection
```http
POST /collections
```

#### Add Media to Collection
```http
POST /collections/{id}/media
```

## 🔔 Webhooks

Media Vault can send webhook notifications for various events:

### Available Events
- `media.uploaded`
- `media.deleted`
- `user.registered`
- `collection.created`

### Webhook Payload

```json
{
  "event": "media.uploaded",
  "data": {
    "id": "media_123",
    "type": "image",
    "user_id": "user_123"
  },
  "timestamp": "2023-05-30T10:00:00Z"
}
```

## ⚠️ Rate Limiting

- **Anonymous**: 60 requests per minute
- **Authenticated**: 1000 requests per minute
- **API Key**: 5000 requests per minute

Headers:
- `X-RateLimit-Limit`: Request limit per time window
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Timestamp when the limit resets

## 📑 Pagination

List endpoints support pagination:

```http
GET /media?page=2&per_page=20
```

Response includes pagination metadata:

```json
{
  "data": [],
  "meta": {
    "pagination": {
      "total": 150,
      "count": 20,
      "per_page": 20,
      "current_page": 2,
      "total_pages": 8,
      "links": {
        "next": "/media?page=3",
        "prev": "/media?page=1"
      }
    }
  }
}
```

## 🌐 WebSocket API

Real-time updates are available via WebSocket:

```javascript
const socket = new WebSocket('wss://api.your-mediavault-instance.com/realtime');

// Authenticate
socket.send(JSON.stringify({
  type: 'auth',
  token: 'your.jwt.token'
}));

// Subscribe to events
socket.send(JSON.stringify({
  type: 'subscribe',
  channel: 'media_updates',
  resource_id: 'media_123'
}));

// Handle messages
socket.onmessage = (event) => {
  console.log('Update:', JSON.parse(event.data));
};
```

## 📚 SDKs

Official SDKs are available for popular languages:

- **JavaScript/TypeScript**: `npm install @mediavault/sdk`
- **Python**: `pip install mediavault-sdk`
- **Go**: `go get github.com/wronai/mediavault-go`

## 📅 API Versioning

The API is versioned in the URL path (e.g., `/api/v1/...`). Breaking changes will result in a new version number.

## 🔒 Security

- All endpoints require HTTPS
- JWT tokens expire after 1 hour
- Password hashing with Argon2id
- Rate limiting to prevent abuse
- CORS restricted to trusted domains

## 📝 Support

For API support, please contact api-support@wron.ai or open an issue on GitHub.