# API Documentation

## Overview

The 3X-UI Modern API is a RESTful API with JWT-based authentication. All endpoints return JSON responses.

**Base URL**: `http://localhost:8080/api/v1`

**Swagger UI**: Available at `http://localhost:8080/docs`

## Authentication

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin"
}
```

**Response**:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

### 2FA Verification

```http
POST /api/v1/auth/2fa/verify
Content-Type: application/json

{
  "username": "admin",
  "code": "123456"
}
```

### Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

## Users

### List Users

```http
GET /api/v1/users?page=1&limit=10
Authorization: Bearer {access_token}
```

### Get User

```http
GET /api/v1/users/{id}
Authorization: Bearer {access_token}
```

### Create User

```http
POST /api/v1/users
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "username": "newuser",
  "email": "user@example.com",
  "password": "securepassword",
  "role": "operator"
}
```

### Update User

```http
PUT /api/v1/users/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "enabled": true,
  "role": "admin"
}
```

### Delete User

```http
DELETE /api/v1/users/{id}
Authorization: Bearer {access_token}
```

## Clients

### List Clients

```http
GET /api/v1/clients?page=1&limit=10&protocol=vmess&enable=true
Authorization: Bearer {access_token}
```

### Get Client

```http
GET /api/v1/clients/{id}
Authorization: Bearer {access_token}
```

### Create Client

```http
POST /api/v1/clients
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "email": "client@example.com",
  "uuid": "generated-uuid",
  "protocol": "vmess",
  "enable": true,
  "traffic_limit": 10737418240,
  "ip_limit": 2
}
```

### Get Client Traffic

```http
GET /api/v1/clients/{id}/traffic
Authorization: Bearer {access_token}
```

## Nodes

### List Nodes

```http
GET /api/v1/nodes?page=1&limit=10
Authorization: Bearer {access_token}
```

### Create Node

```http
POST /api/v1/nodes
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Node 1",
  "address": "node1.example.com",
  "port": 443,
  "status": "online"
}
```

## Metrics

### Get Summary

```http
GET /api/v1/metrics/summary
Authorization: Bearer {access_token}
```

**Response**:
```json
{
  "users": {
    "total": 5
  },
  "clients": {
    "total": 100,
    "active": 85
  },
  "nodes": {
    "total": 3,
    "online": 2
  }
}
```

## WebSocket

### Monitor Endpoint

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/monitor');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Update:', data);
};
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

**HTTP Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Too Many Requests
- `500` - Internal Server Error

## Rate Limiting

- Default: 100 requests per minute per IP
- Authenticated: 1000 requests per minute per user

## Pagination

All list endpoints support pagination:

```http
GET /api/v1/users?page=1&limit=10
```

Response includes:
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 10
}
```
