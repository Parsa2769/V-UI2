# V-UI - Advanced Xray Multi-Protocol Panel

<div align="center">

![Version](https://img.shields.io/badge/version-2.0.0-blue)
![License](https://img.shields.io/badge/license-GPL--3.0-green)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)
![Xray](https://img.shields.io/badge/Xray-Latest-9cf)

</div>

[English](#english) | [فارسی](#فارسی)

---

## English

### 🚀 Overview

**V-UI** is a powerful, modern web-based control panel for managing Xray-core with complete protocol support. Built with Go and React, it provides enterprise-grade features including multi-protocol support, advanced traffic management, SSL certificate automation, Telegram bot integration, and much more.

V-UI offers a completely redesigned architecture with enhanced security, modern UI, and production-ready infrastructure.

### 🚀 Quick Install

Install with one command:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

This will automatically:
- ✅ Install Docker & Docker Compose
- ✅ Clone the repository
- ✅ Generate secure secrets
- ✅ Start all services

**Access the panel:**
- Web Panel: `http://YOUR-SERVER-IP`
- API Docs: `http://YOUR-SERVER-IP/docs`

**Default Login:**
- Username: `admin`
- Password: `admin`

⚠️ **IMPORTANT**: Change the default password immediately after first login!

### ⚠️ Legal Notice

**FOR PERSONAL USE ONLY**. This tool is designed for legitimate privacy and security purposes. Any illegal use is strictly prohibited and against the project's intentions. Users are solely responsible for compliance with local laws and regulations.

### ✨ Key Features

#### Core Functionality (Maintained from Original)
- ✅ Full Xray protocol support: VMESS, VLESS, Trojan, Shadowsocks, WireGuard, Tunnel, Mixed, HTTP
- ✅ User management with traffic, day, and IP limitations
- ✅ Multi-node support with centralized management
- ✅ Configuration templates and QR code generation
- ✅ Comprehensive traffic monitoring and logging

#### New Features (v2.0)
- 🎨 **Modern React UI**: Beautiful, responsive interface with dark/light themes
- 🔐 **Enhanced Security**: JWT authentication, 2FA (TOTP), RBAC with admin/operator/viewer roles
- 📊 **Real-time Monitoring**: WebSocket-based live updates for traffic and node status
- 📈 **Advanced Analytics**: Prometheus metrics, detailed reports, and visualizations
- 🔄 **Import/Export**: CSV/JSON support for bulk user management
- 🎯 **Config Templates**: Pre-defined templates for quick deployment
- 🔔 **Alert System**: Slack/Telegram/Email notifications for quota and errors
- 📚 **API Documentation**: Complete OpenAPI/Swagger specs
- 🧪 **Testing**: 60%+ code coverage with unit and E2E tests
- 🚢 **DevOps Ready**: Docker, Kubernetes, CI/CD with GitHub Actions

### 📋 Prerequisites

- **OS**: Linux (Debian/Ubuntu recommended), macOS, Windows (WSL2)
- **Docker**: 20.10+ and Docker Compose 2.0+
- **Go**: 1.21+ (for development)
- **Node.js**: 18+ (for frontend development)
- **Database**: PostgreSQL 14+ (or SQLite for testing)

### 🚀 Quick Start

#### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Copy environment file
cp .env.example .env

# Edit .env with your settings (database, secrets, etc.)
nano .env

# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f
```

Access the panel at `http://localhost:8080`

Default credentials:
- Username: `admin`
- Password: `admin` (⚠️ Change immediately!)

#### Option 2: Manual Installation

```bash
# Install dependencies
./scripts/install.sh

# Run database migrations
./scripts/migrate.sh up

# Build backend
cd backend
go build -o ../bin/v-ui ./cmd/server

# Build frontend
cd ../frontend
npm install
npm run build

# Start server
cd ..
./bin/v-ui --config config.yaml
```

### 🏗️ Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Browser   │─────▶│  React SPA   │─────▶│  Go Backend │
└─────────────┘      │  (TypeScript)│      │   (Xray)    │
                     └──────────────┘      └─────────────┘
                            │                      │
                            │                      │
                     ┌──────▼──────┐      ┌────────▼──────┐
                     │  TailwindCSS│      │  PostgreSQL   │
                     │  shadcn/ui  │      │   /SQLite     │
                     └─────────────┘      └───────────────┘
```

### 📁 Project Structure

```
V-UI2/
├── backend/                 # Go backend
│   ├── cmd/                # Entry points
│   ├── internal/           # Private application code
│   │   ├── api/           # API handlers
│   │   ├── auth/          # Authentication & authorization
│   │   ├── config/        # Configuration management
│   │   ├── db/            # Database layer
│   │   ├── middleware/    # HTTP middleware
│   │   ├── models/        # Data models
│   │   ├── service/       # Business logic
│   │   └── xray/          # Xray integration
│   └── pkg/               # Public libraries
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/    # Reusable components
│   │   ├── pages/         # Page components
│   │   ├── hooks/         # Custom hooks
│   │   ├── utils/         # Utilities
│   │   ├── api/           # API client
│   │   └── types/         # TypeScript types
│   ├── public/
│   └── package.json
├── migrations/             # Database migrations
├── docs/                   # Documentation
├── scripts/                # Deployment scripts
├── k8s/                    # Kubernetes manifests
├── .github/                # CI/CD workflows
├── docker-compose.yml
├── Dockerfile
└── README.md
```

### 🔧 Configuration

Configuration is managed via `config.yaml` or environment variables:

```yaml
server:
  port: 8080
  mode: production

database:
  type: postgres
  host: localhost
  port: 5432
  name: x3ui
  user: x3ui
  password: ${DB_PASSWORD}

auth:
  jwt_secret: ${JWT_SECRET}
  token_expiry: 24h
  refresh_expiry: 168h
  enable_2fa: true

xray:
  config_path: /etc/xray/config.json
  api_port: 10085

monitoring:
  enable_metrics: true
  enable_websocket: true
```

### 🧪 Testing

```bash
# Backend tests
cd backend
go test ./... -v -cover

# Frontend tests
cd frontend
npm run test

# E2E tests
npm run test:e2e

# Run all tests
./scripts/test.sh
```

### 📚 API Documentation

API documentation is available at `/docs` when running the server:
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI Spec: `http://localhost:8080/api/v1/openapi.json`

Key endpoints:
- `POST /api/v1/auth/login` - User login (returns JWT)
- `POST /api/v1/auth/2fa/verify` - 2FA verification
- `GET /api/v1/users` - List users (requires auth)
- `POST /api/v1/users` - Create user
- `GET /api/v1/nodes` - List nodes
- `GET /api/v1/metrics` - Prometheus metrics
- `GET /ws/monitor` - WebSocket live updates

### 🔒 Security Best Practices

1. **Change default credentials immediately**
2. **Use strong JWT secrets** (32+ random characters)
3. **Enable 2FA** for all admin accounts
4. **Use HTTPS** in production (configure reverse proxy)
5. **Regular backups** of database and configuration
6. **Keep dependencies updated** (`dependabot` enabled)
7. **Review audit logs** regularly
8. **Limit API access** with rate limiting and RBAC

### 🔄 Data Import

```bash
# V-UI supports importing configurations
# Use the backup/restore feature in settings
# Or import via API endpoints
```

### 🔄 Previous Installation Data

```bash
# Import existing data
./scripts/backup.sh

# Use the backup/restore feature
# Or import via API endpoints

# Verify data
./scripts/verify-migration.sh
```

### 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### 📝 License

GPL-3.0 License - see [LICENSE](LICENSE) file.

### 🙏 Acknowledgments

- **V-UI** - Modern Xray management panel with clean architecture
- Inspired by: [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui)
- Core engine: [Xray-core](https://github.com/XTLS/Xray-core)
- Open source community contributors

---

## فارسی

### 🚀 معرفی

**V-UI** یک پنل کنترل وب مدرن و قدرتمند برای مدیریت پروتکل‌های Xray با معماری تمیز، رابط کاربری پیشرفته و قابلیت‌های سازمانی است. این پروژه با Go و React ساخته شده و زیرساخت آماده برای تولید دارد.

### ⚠️ اطلاعیه قانونی

**فقط برای استفاده شخصی**. این ابزار برای اهداف مشروع حریم خصوصی و امنیت طراحی شده است. هرگونه استفاده غیرقانونی اکیداً ممنوع و بر خلاف اهداف پروژه است. کاربران مسئولیت انطباق با قوانین و مقررات محلی را بر عهده دارند.

### ✨ ویژگی‌های کلیدی

#### عملکرد اصلی (حفظ شده از نسخه اصلی)
- ✅ پشتیبانی کامل از پروتکل‌های Xray: VMESS، VLESS، Trojan، Shadowsocks، WireGuard، Tunnel، Mixed، HTTP
- ✅ مدیریت کاربران با محدودیت ترافیک، روز و IP
- ✅ پشتیبانی چند نود با مدیریت متمرکز
- ✅ قالب‌های پیکربندی و تولید کد QR
- ✅ مانیتورینگ جامع ترافیک و لاگ‌گذاری

#### ویژگی‌های جدید (نسخه 2.0)
- 🎨 **رابط کاربری مدرن React**: رابط زیبا و واکنش‌گرا با تم تاریک/روشن
- 🔐 **امنیت پیشرفته**: احراز هویت JWT، 2FA (TOTP)، RBAC با نقش‌های admin/operator/viewer
- 📊 **مانیتورینگ بلادرنگ**: به‌روزرسانی زنده مبتنی بر WebSocket برای ترافیک و وضعیت نود
- 📈 **تحلیل پیشرفته**: معیارهای Prometheus، گزارش‌های دقیق و بصری‌سازی
- 🔄 **Import/Export**: پشتیبانی از CSV/JSON برای مدیریت انبوه کاربران
- 🎯 **قالب‌های پیکربندی**: قالب‌های از پیش تعریف شده برای استقرار سریع
- 🔔 **سیستم هشدار**: اعلان‌های Slack/Telegram/Email برای سهمیه و خطاها
- 📚 **مستندات API**: مشخصات کامل OpenAPI/Swagger
- 🧪 **تست**: پوشش کد بیش از 60% با تست‌های واحد و E2E
- 🚢 **آماده DevOps**: Docker، Kubernetes، CI/CD با GitHub Actions

### 📋 پیش‌نیازها

- **سیستم عامل**: Linux (Debian/Ubuntu توصیه شده)، macOS، Windows (WSL2)
- **Docker**: نسخه 20.10+ و Docker Compose 2.0+
- **Go**: 1.21+ (برای توسعه)
- **Node.js**: 18+ (برای توسعه فرانت‌اند)
- **پایگاه داده**: PostgreSQL 14+ (یا SQLite برای تست)

### 🚀 شروع سریع

#### گزینه 1: Docker Compose (توصیه شده)

```bash
# کلون کردن مخزن
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# کپی کردن فایل محیطی
cp .env.example .env

# ویرایش .env با تنظیمات خود
nano .env

# شروع تمام سرویس‌ها
docker-compose up -d

# مشاهده لاگ‌ها
docker-compose logs -f
```

دسترسی به پنل در آدرس `http://localhost:8080`

اطلاعات ورود پیش‌فرض:
- نام کاربری: `admin`
- رمز عبور: `admin` (⚠️ فوراً تغییر دهید!)

### 🏗️ معماری

پروژه از معماری مدرن سه‌لایه استفاده می‌کند:
- **Frontend**: React + TypeScript + TailwindCSS + shadcn/ui
- **Backend**: Go با معماری Clean Architecture
- **Database**: PostgreSQL/SQLite با مهاجرت‌های نسخه‌بندی شده

### 🧪 تست

```bash
# تست‌های بک‌اند
cd backend
go test ./... -v -cover

# تست‌های فرانت‌اند
cd frontend
npm run test

# تست‌های E2E
npm run test:e2e
```

### 📚 مستندات API

مستندات API در آدرس `/docs` در دسترس است:
- Swagger UI: `http://localhost:8080/docs`
- مشخصات OpenAPI: `http://localhost:8080/api/v1/openapi.json`

### 🔒 بهترین شیوه‌های امنیتی

1. **فوراً اطلاعات ورود پیش‌فرض را تغییر دهید**
2. **از رمزهای قوی JWT استفاده کنید** (32+ کاراکتر تصادفی)
3. **2FA را فعال کنید** برای تمام حساب‌های مدیر
4. **از HTTPS استفاده کنید** در محیط تولید
5. **پشتیبان‌گیری منظم** از پایگاه داده و پیکربندی
6. **وابستگی‌ها را به‌روز نگه دارید**
7. **لاگ‌های ممیزی را به طور منظم بررسی کنید**
8. **دسترسی API را محدود کنید** با rate limiting و RBAC

### 🔄 وارد کردن داده‌ها

```bash
# V-UI از import کانفیگ پشتیبانی می‌کند
# از قابلیت backup/restore در تنظیمات استفاده کنید
# یا از طریق API endpoint ها import کنید
```

### 🔄 داده‌های نصب قبلی

```bash
# وارد کردن داده‌های موجود
./scripts/backup.sh

# از قابلیت backup/restore استفاده کنید
# یا از طریق API endpoint ها import کنید

# تأیید داده‌ها
./scripts/verify-migration.sh
```

### 🤝 مشارکت

برای راهنمای مشارکت، [CONTRIBUTING.md](CONTRIBUTING.md) را مشاهده کنید.

### 📝 مجوز

مجوز GPL-3.0 - فایل [LICENSE](LICENSE) را مشاهده کنید.

### 🙏 قدردانی

- **V-UI** - پنل مدیریت مدرن Xray با معماری تمیز
- الهام گرفته از: [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui)
- هسته اصلی: [Xray-core](https://github.com/XTLS/Xray-core)
- مشارکت‌کنندگان open source

---

## 📞 Support

- 📖 [Documentation](./docs)
- 🐛 [Issue Tracker](https://github.com/Parsa2769/V-UI2/issues)
- 💬 [Discussions](https://github.com/Parsa2769/V-UI2/discussions)

**Remember**: Use responsibly and legally. This project is for legitimate privacy and security purposes only.
