# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-10-14

### Added

#### Frontend
- ✨ Complete React + TypeScript rewrite with modern UI/UX
- 🎨 TailwindCSS + shadcn/ui components for beautiful, accessible design
- 🌓 Dark/Light theme support with smooth transitions
- 📱 Fully responsive design for mobile, tablet, and desktop
- 📊 Real-time dashboard with KPI cards and live charts
- 🔄 WebSocket integration for live traffic and node status updates
- 📈 Advanced analytics and visualization with recharts
- 🔍 Advanced filtering, sorting, and pagination for all data tables
- 📥 CSV/JSON export and import functionality
- 🎯 QR code generator for configuration sharing
- 🔔 Toast notifications and alert system

#### Backend
- 🏗️ Modular architecture with clean separation of concerns
- 🔐 JWT-based authentication with access and refresh tokens
- 🔒 2FA (TOTP) implementation with QR code setup
- 👥 RBAC (Role-Based Access Control): admin, operator, viewer roles
- 📡 WebSocket server for real-time updates
- 📊 Prometheus metrics endpoint for monitoring
- 🔄 Background job processing with worker queues
- 📝 Structured logging with zap
- 🔒 Rate limiting and input validation
- 💾 Database migration system with golang-migrate
- 🗄️ Support for both SQLite (dev) and PostgreSQL (prod)
- 🔍 Audit logging for security-sensitive actions

#### API
- 📚 Complete OpenAPI/Swagger documentation
- 🌐 RESTful API with consistent error handling
- 🔌 WebSocket endpoint for live monitoring
- 📄 Pagination, filtering, and sorting support
- 🔐 Per-endpoint authentication and authorization
- 📊 Comprehensive metrics and health check endpoints

#### DevOps
- 🐳 Multi-stage Dockerfile for optimized builds
- 🚀 Docker Compose for easy local development
- ☸️ Kubernetes manifests and Helm charts
- 🔄 GitHub Actions CI/CD pipelines
- 🧪 Automated testing in CI (unit + E2E)
- 🔍 Linting and code quality checks
- 📦 Automated dependency updates with Dependabot
- 🏷️ Automated releases with Docker image publishing

#### Testing
- 🧪 Unit tests for backend services (60%+ coverage)
- 🎭 E2E tests with Playwright
- 🔬 Integration tests with test database
- 📊 Test coverage reporting
- ✅ Table-driven tests for Go code
- 🎯 Component tests for React with Testing Library

#### Documentation
- 📖 Comprehensive README in English and Farsi
- 📚 API documentation with Swagger UI
- 🏗️ Architecture diagrams and design docs
- 🔄 Migration guide from original 3x-ui
- 🛠️ Deployment guides (Docker, K8s, bare metal)
- 🔒 Security best practices guide
- 🤝 Contributing guidelines
- 📝 Code of conduct

### Changed
- 🔧 Improved configuration management with environment variables
- ⚡ Enhanced performance with connection pooling and caching
- 🔒 Strengthened security with bcrypt password hashing
- 📦 Updated all dependencies to latest stable versions
- 🎨 Modernized UI/UX based on contemporary design principles

### Fixed
- 🐛 Various bug fixes from original project
- 🔒 Security vulnerabilities in dependencies
- 💾 Database connection leaks
- 🔄 Race conditions in concurrent operations

### Maintained from Original
- ✅ All Xray protocol support (VMESS, VLESS, Trojan, Shadowsocks, WireGuard)
- ✅ User management with traffic/day/IP limits
- ✅ Multi-node support
- ✅ Configuration templates
- ✅ Traffic monitoring and logging
- ✅ Backward compatibility with existing configs

## [1.0.0] - Original 3x-ui

Base functionality from MHSanaei/3x-ui project.

---

## Upgrade Guide

### From 1.x to 2.0

**⚠️ Important**: Back up your data before upgrading!

```bash
# Backup existing data
./scripts/backup.sh

# Stop old version
docker-compose down

# Pull new version
git pull origin main

# Run migrations
./scripts/migrate.sh up

# Start new version
docker-compose up -d

# Verify
./scripts/verify-migration.sh
```

### Breaking Changes in 2.0

1. **Database Schema**: New tables for RBAC, audit logs, and 2FA
2. **API**: New authentication endpoints (JWT instead of session)
3. **Configuration**: Updated config format (migrate with migration tool)
4. **Environment Variables**: New required variables (see .env.example)

### New Required Environment Variables

```env
JWT_SECRET=<your-32-char-secret>
JWT_REFRESH_SECRET=<your-32-char-secret>
ENABLE_2FA=true
DATABASE_URL=postgresql://user:pass@host:5432/dbname
```

See [MIGRATION.md](docs/MIGRATION.md) for detailed upgrade instructions.
