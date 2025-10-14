# Quick Start Guide

Get 3X-UI Modern running in under 5 minutes!

## Prerequisites

- **Docker** 20.10+ and **Docker Compose** 2.0+
- **2GB RAM** minimum
- **10GB disk space**

## Installation Steps

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/3x-ui-modern.git
cd 3x-ui-modern

# Run installation script (Linux/macOS)
chmod +x scripts/install.sh
./scripts/install.sh

# Or manually copy environment file
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` file and set secure values:

```bash
# Generate secure secrets
JWT_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)
DATABASE_PASSWORD=$(openssl rand -base64 24)

# Update .env file with these values
nano .env
```

**Critical settings to update:**
- `JWT_SECRET` - Must be at least 32 characters
- `JWT_REFRESH_SECRET` - Must be at least 32 characters  
- `DATABASE_PASSWORD` - Strong database password
- `DEFAULT_ADMIN_PASSWORD` - Change from default

### 3. Start Services

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### 4. Access the Panel

Open your browser and navigate to:

**URL**: `http://localhost:8080`

**Default Credentials**:
- Username: `admin`
- Password: `admin`

**⚠️ IMPORTANT**: Change the password immediately after first login!

### 5. Verify Installation

```bash
# Check health endpoint
curl http://localhost:8080/health

# Check API
curl http://localhost:8080/api/v1/metrics/summary
```

## Using Make Commands

```bash
# View all commands
make help

# Start services
make start

# Stop services
make stop

# View logs
make logs

# Create backup
make backup

# Run tests
make test
```

## Next Steps

1. **Change Admin Password**
   - Login to the panel
   - Go to Settings → Profile
   - Change password

2. **Enable 2FA** (Highly Recommended)
   - Settings → Security
   - Enable Two-Factor Authentication
   - Scan QR code with authenticator app

3. **Create Additional Users**
   - Navigate to Users page
   - Click "Create User"
   - Assign appropriate role (admin/operator/viewer)

4. **Add Nodes**
   - Go to Nodes page
   - Add your Xray server nodes

5. **Create Clients**
   - Navigate to Clients page
   - Create new clients with desired protocols

## API Access

### Get API Token

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Response includes access_token
```

### Use Token

```bash
# List users
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### API Documentation

Visit `http://localhost:8080/docs` for interactive Swagger documentation.

## Monitoring (Optional)

Enable monitoring stack with Prometheus and Grafana:

```bash
# Start with monitoring
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

# Restart services
docker-compose restart

# Rebuild
docker-compose build --no-cache
docker-compose up -d
```

### Can't connect to database

```bash
# Check database is running
docker-compose ps postgres

# View database logs
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d
```

### Frontend not loading

```bash
# Check backend logs
docker-compose logs backend

# Rebuild frontend
docker-compose build frontend
docker-compose up -d
```

### Port already in use

```bash
# Change port in .env
SERVER_PORT=8081

# Or stop conflicting service
sudo lsof -ti:8080 | xargs kill -9
```

## Security Checklist

- [ ] Changed default admin password
- [ ] Set strong JWT secrets (32+ characters)
- [ ] Enabled 2FA for admin accounts
- [ ] Configured HTTPS (for production)
- [ ] Reviewed `.env` file for sensitive data
- [ ] Enabled firewall rules
- [ ] Regular backups configured
- [ ] Updated default database password

## Common Commands

```bash
# Create backup
./scripts/backup.sh

# Update to latest version
./scripts/update.sh

# View real-time logs
docker-compose logs -f backend

# Execute command in container
docker-compose exec backend sh

# Database shell
docker-compose exec postgres psql -U x3ui

# Stop all services
docker-compose down

# Remove all data (⚠️ destructive)
docker-compose down -v
```

## Development Mode

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Run tests
npm test

# Run E2E tests
npm run test:e2e
```

## Getting Help

- **Documentation**: `./docs/` directory
- **API Docs**: http://localhost:8080/docs
- **Issues**: https://github.com/yourusername/3x-ui-modern/issues
- **Discussions**: https://github.com/yourusername/3x-ui-modern/discussions

## Production Deployment

For production deployment, see:
- [DEPLOYMENT.md](./docs/DEPLOYMENT.md) - Comprehensive deployment guide
- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - System architecture
- [API.md](./docs/API.md) - API documentation

## License

GPL-3.0 - See [LICENSE](./LICENSE) file.

## Legal Notice

**FOR PERSONAL USE ONLY**. This tool is designed for legitimate privacy and security purposes. Any illegal use is strictly prohibited. Users are solely responsible for compliance with local laws.
