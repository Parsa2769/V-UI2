# Deployment Guide

## Docker Compose (Recommended)

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum
- 10GB disk space

### Quick Start

```bash
# Clone repository
git clone https://github.com/yourusername/3x-ui-modern.git
cd 3x-ui-modern

# Run installation script
chmod +x scripts/install.sh
./scripts/install.sh

# Edit configuration
nano .env

# Start services
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Access

- **Web Panel**: http://localhost:8080
- **API Docs**: http://localhost:8080/docs
- **Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090 (if monitoring profile enabled)
- **Grafana**: http://localhost:3000 (if monitoring profile enabled)

## Manual Installation

### Backend

```bash
# Install Go 1.21+
cd backend

# Install dependencies
go mod download

# Build
go build -o ../bin/3x-ui-modern ./cmd/server

# Run migrations
./bin/3x-ui-modern migrate

# Start server
./bin/3x-ui-modern
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Development
npm run dev

# Production build
npm run build
```

## Kubernetes

### Using Kubectl

```bash
# Create namespace
kubectl create namespace 3x-ui-modern

# Apply manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -n 3x-ui-modern

# Get service URL
kubectl get svc -n 3x-ui-modern
```

### Using Helm

```bash
# Add repository
helm repo add 3x-ui-modern https://yourusername.github.io/3x-ui-modern

# Install
helm install my-3x-ui 3x-ui-modern/3x-ui-modern \
  --set database.password=SecurePassword \
  --set auth.jwtSecret=Your32CharSecret

# Upgrade
helm upgrade my-3x-ui 3x-ui-modern/3x-ui-modern

# Uninstall
helm uninstall my-3x-ui
```

## Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/ssl/certs/yourdomain.crt;
    ssl_certificate_key /etc/ssl/private/yourdomain.key;
    
    # Frontend
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

## Environment Variables

Key variables to configure:

```bash
# Security (REQUIRED)
JWT_SECRET=your-32-char-secret
JWT_REFRESH_SECRET=your-32-char-refresh-secret
DATABASE_PASSWORD=secure-password

# Server
SERVER_PORT=8080
SERVER_MODE=production

# Database
DATABASE_TYPE=postgres
DATABASE_HOST=postgres
DATABASE_NAME=x3ui

# Features
ENABLE_2FA=true
ENABLE_METRICS=true
ENABLE_WEBSOCKET=true
```

## SSL/TLS Configuration

### Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com

# Auto-renewal (cron)
0 0 * * * certbot renew --quiet
```

## Backup and Restore

### Backup

```bash
# Automated backup
./scripts/backup.sh

# Manual database backup
docker-compose exec postgres pg_dump -U x3ui x3ui > backup.sql
```

### Restore

```bash
# Stop services
docker-compose down

# Restore database
docker-compose up -d postgres
cat backup.sql | docker-compose exec -T postgres psql -U x3ui x3ui

# Restore data
tar -xzf backup_TIMESTAMP.tar.gz

# Start services
docker-compose up -d
```

## Monitoring

### Enable Monitoring Stack

```bash
# Start with monitoring profile
docker-compose --profile monitoring up -d

# Access Grafana
http://localhost:3000
# Default: admin/admin

# Access Prometheus
http://localhost:9090
```

## Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs

# Check disk space
df -h

# Check permissions
ls -la data/
```

### Database connection failed

```bash
# Check database status
docker-compose ps postgres

# Check connection
docker-compose exec backend nc -zv postgres 5432

# Reset database
docker-compose down -v
docker-compose up -d
```

### High memory usage

```bash
# Check resource usage
docker stats

# Limit resources in docker-compose.yml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
```

## Security Hardening

1. **Change default credentials immediately**
2. **Use strong JWT secrets** (32+ random characters)
3. **Enable HTTPS** with valid certificates
4. **Enable 2FA** for all admin accounts
5. **Configure firewall** to limit access
6. **Regular backups** (automated recommended)
7. **Keep dependencies updated**
8. **Monitor audit logs** regularly
9. **Use rate limiting** on public endpoints
10. **Restrict database access** to internal network only

## Performance Tuning

### Database

```bash
# Increase connections
DATABASE_MAX_OPEN_CONNS=50
DATABASE_MAX_IDLE_CONNS=10
```

### Backend

```bash
# Enable caching
# Add Redis for session storage
```

### Frontend

```bash
# Use CDN for static assets
# Enable gzip compression in nginx
```
