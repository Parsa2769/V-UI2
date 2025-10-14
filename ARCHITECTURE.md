# Architecture Documentation

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Browser  │  │  Mobile  │  │   CLI    │  │   API    │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        └─────────────┴─────────────┴─────────────┘
                        │
        ┌───────────────▼────────────────┐
        │      Nginx Reverse Proxy       │
        │         (SSL/TLS)              │
        └───────────────┬────────────────┘
                        │
        ┌───────────────▼────────────────┐
        │        React Frontend          │
        │     (Vite + TypeScript)        │
        └───────────────┬────────────────┘
                        │
        ┌───────────────▼────────────────┐
        │        Go Backend API          │
        │    ┌──────────────────────┐    │
        │    │   HTTP Handlers      │    │
        │    ├──────────────────────┤    │
        │    │   Middleware         │    │
        │    ├──────────────────────┤    │
        │    │   Business Logic     │    │
        │    ├──────────────────────┤    │
        │    │   Data Access        │    │
        │    └──────────────────────┘    │
        └───────────────┬────────────────┘
                        │
        ┌───────────────▼────────────────┐
        │       PostgreSQL/SQLite        │
        │      (Persistent Storage)      │
        └────────────────────────────────┘
                        │
        ┌───────────────▼────────────────┐
        │          Xray Core             │
        │      (Protocol Handler)        │
        └────────────────────────────────┘
```

## Backend Architecture

### Layered Architecture

```
cmd/
  └── server/        # Application entry point
internal/
  ├── api/          # HTTP handlers & routing
  ├── auth/         # Authentication & authorization
  ├── config/       # Configuration management
  ├── db/           # Database layer
  ├── middleware/   # HTTP middleware
  ├── models/       # Data models
  ├── service/      # Business logic
  ├── websocket/    # WebSocket server
  └── xray/         # Xray integration
pkg/                # Public libraries
```

### Request Flow

```
HTTP Request
    │
    ▼
Middleware Chain
  ├── Logger
  ├── Recovery
  ├── CORS
  ├── RequestID
  ├── RateLimit
  └── Auth (if protected)
    │
    ▼
Router (Gin)
    │
    ▼
Handler (API Layer)
    │
    ▼
Service Layer (Business Logic)
    │
    ▼
Repository/DB Layer
    │
    ▼
Database (GORM)
```

## Frontend Architecture

### Component Structure

```
src/
├── components/
│   ├── layout/      # Layout components
│   ├── common/      # Reusable components
│   └── features/    # Feature-specific components
├── pages/           # Page components
├── hooks/           # Custom React hooks
├── stores/          # State management (Zustand)
├── lib/             # Utilities & API client
└── types/           # TypeScript types
```

### State Management

- **Zustand**: Global state (auth, theme)
- **React Query**: Server state & caching
- **Local State**: Component-specific state

### Data Flow

```
Component
    │
    ▼
React Query (useQuery/useMutation)
    │
    ▼
API Client (Axios)
    │
    ▼
Backend API
    │
    ▼
Component Update (React)
```

## Authentication Flow

### Login Flow

```
1. User submits credentials
2. Backend validates credentials
3. If 2FA enabled → request code
4. Backend generates JWT tokens
5. Tokens stored in localStorage
6. Access token used in API requests
7. Refresh token used to renew access token
```

### JWT Structure

```
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload:
{
  "user_id": "uuid",
  "username": "admin",
  "role": "admin",
  "exp": 1234567890
}
```

## Database Schema

### Core Tables

```sql
users
  - id (UUID, PK)
  - username (unique)
  - email (unique)
  - password_hash
  - role (admin|operator|viewer)
  - enabled
  - two_factor_secret
  - two_factor_enabled
  - last_login
  - created_at
  - updated_at

clients
  - id (UUID, PK)
  - email (unique)
  - uuid (unique)
  - protocol
  - enable
  - expiry_time
  - traffic_limit
  - upload_bytes
  - download_bytes
  - ip_limit
  - node_id (FK)
  - created_by (FK)
  - created_at
  - updated_at

nodes
  - id (UUID, PK)
  - name (unique)
  - address
  - port
  - status
  - last_seen
  - upload_bytes
  - download_bytes
  - metadata (JSONB)
  - created_at
  - updated_at

traffic_logs
  - id (UUID, PK)
  - client_id (FK)
  - node_id (FK)
  - upload_bytes
  - download_bytes
  - timestamp

audit_logs
  - id (UUID, PK)
  - actor_id (FK)
  - action
  - target
  - detail
  - ip_address
  - user_agent
  - created_at

refresh_tokens
  - id (UUID, PK)
  - user_id (FK)
  - token (unique)
  - expires_at
  - created_at
  - revoked_at
```

## Security Architecture

### Defense in Depth

1. **Transport Layer**: HTTPS/TLS
2. **Network Layer**: Firewall, VPN
3. **Application Layer**: 
   - JWT authentication
   - RBAC authorization
   - Rate limiting
   - Input validation
4. **Data Layer**:
   - Encrypted passwords (bcrypt)
   - Encrypted secrets
   - Audit logging

### RBAC Model

```
Roles:
  Admin    → Full access
  Operator → CRUD clients, nodes
  Viewer   → Read-only access

Resources:
  - Users
  - Clients
  - Nodes
  - Traffic
  - Audit Logs
  - Settings
```

## Monitoring & Observability

### Metrics (Prometheus)

- HTTP request count & duration
- Active connections
- Traffic bytes
- Active clients
- Nodes online

### Logging (Structured)

```json
{
  "level": "info",
  "time": "2024-01-01T00:00:00Z",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "12ms",
  "ip": "192.168.1.1"
}
```

### Health Checks

- `/health` - Basic health status
- `/ready` - Readiness probe (DB check)

## Scalability Considerations

### Horizontal Scaling

- **Stateless Backend**: Can scale horizontally
- **Database**: Use connection pooling
- **Sessions**: Store in Redis (if needed)
- **Load Balancer**: Nginx or cloud LB

### Vertical Scaling

- Increase resources per container
- Optimize database queries
- Use caching (Redis)

## Performance Optimization

### Backend

- Connection pooling
- Query optimization
- Caching (Redis)
- Goroutines for async tasks

### Frontend

- Code splitting
- Lazy loading
- Image optimization
- CDN for static assets

### Database

- Indexes on frequently queried fields
- Query optimization
- Regular VACUUM (PostgreSQL)
- Connection pooling

## Deployment Architecture

### Docker Compose

```
Services:
  - postgres (database)
  - backend (Go API)
  - frontend (React SPA)
  - nginx (reverse proxy)
  - prometheus (metrics)
  - grafana (dashboards)
```

### Kubernetes

```
Deployments:
  - backend (replicas: 3)
  - frontend (replicas: 2)
  
StatefulSets:
  - postgres (replicas: 1)
  
Services:
  - backend-svc (ClusterIP)
  - postgres-svc (ClusterIP)
  - nginx-svc (LoadBalancer)
  
ConfigMaps:
  - app-config
  
Secrets:
  - jwt-secrets
  - db-credentials
```

## Technology Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin
- **ORM**: GORM
- **Auth**: golang-jwt
- **Logging**: Zap
- **Metrics**: Prometheus client
- **WebSocket**: gorilla/websocket

### Frontend
- **Framework**: React 18
- **Language**: TypeScript 5
- **Build Tool**: Vite
- **Styling**: TailwindCSS
- **State**: Zustand, React Query
- **HTTP Client**: Axios
- **Router**: React Router

### Infrastructure
- **Container**: Docker
- **Orchestration**: Kubernetes (optional)
- **Database**: PostgreSQL 16
- **Proxy**: Nginx
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions
