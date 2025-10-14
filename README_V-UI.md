# V-UI - Advanced Xray Control Panel

<div align="center">

![V-UI Logo](https://via.placeholder.com/150x150/4A90E2/FFFFFF?text=V-UI)

[![Version](https://img.shields.io/badge/version-2.0.0-blue)](https://github.com/Parsa2769/V-UI2)
[![License](https://img.shields.io/badge/license-GPL--3.0-green)](./LICENSE)
[![CI](https://github.com/Parsa2769/V-UI2/workflows/CI/badge.svg)](https://github.com/Parsa2769/V-UI2/actions)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

**پنل مدیریت پیشرفته Xray با رابط کاربری مدرن**

[English](#english) | [فارسی](#persian)

</div>

---

## 🌟 معرفی

**V-UI** یک پنل کنترل وب پیشرفته برای مدیریت Xray-core با پشتیبانی کامل از تمام پروتکل‌ها است. این پنل با معماری مدرن Go و React ساخته شده و امکانات enterprise-grade از جمله:

- ✅ پشتیبانی کامل از **تمام پروتکل‌های Xray**
- ✅ **مدیریت چندکاربره** با محدودیت‌های ترافیک و IP
- ✅ **رابط کاربری مدرن** با تم تاریک/روشن
- ✅ **احراز هویت پیشرفته** (JWT + 2FA + RBAC)
- ✅ **مدیریت گواهینامه SSL** (Let's Encrypt + خودکار)
- ✅ **ربات تلگرام** برای مدیریت و نظارت
- ✅ **Subscription System** کامل
- ✅ **آمارگیری real-time** با WebSocket
- ✅ **API مستند** با Swagger
- ✅ **Docker** و **CI/CD** آماده

---

## 🚀 نصب سریع

نصب با یک دستور:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

این دستور به صورت خودکار:
- ✅ Docker و Docker Compose را نصب می‌کند
- ✅ پروژه را کلون می‌کند
- ✅ Secret های امن تولید می‌کند
- ✅ تمام سرویس‌ها را اجرا می‌کند

### دسترسی به پنل

پس از نصب:
- **پنل وب**: `http://YOUR-SERVER-IP`
- **مستندات API**: `http://YOUR-SERVER-IP/docs`

**اطلاعات ورود پیش‌فرض:**
- نام کاربری: `admin`
- رمز عبور: `admin`

⚠️ **مهم**: بلافاصله پس از اولین ورود رمز عبور را تغییر دهید!

---

## 📋 فهرست مطالب

- [ویژگی‌ها](#-ویژگیها)
- [پروتکل‌های پشتیبانی شده](#-پروتکلهای-پشتیبانی-شده)
- [نصب](#-نصب)
- [پیکربندی](#-پیکربندی)
- [استفاده](#-استفاده)
- [API](#-api)
- [مستندات](#-مستندات)
- [توسعه](#-توسعه)
- [مجوز](#-مجوز)

---

## ✨ ویژگی‌ها

### 🔐 امنیت پیشرفته
- **JWT Authentication** با Access و Refresh Token
- **2FA (TOTP)** با QR Code
- **RBAC** با 3 نقش (Admin, Operator, Viewer)
- **Bcrypt** برای هش کردن رمز عبور (cost 12)
- **Rate Limiting** برای API
- **Audit Logging** برای تمام عملیات
- **Input Validation** کامل

### 🌐 پروتکل‌ها
- **VMess** - با تمام تنظیمات
- **VLESS** - با پشتیبانی XTLS (Vision, Direct, Reality)
- **Trojan** - با TLS/Reality
- **Shadowsocks** - تمام روش‌های رمزنگاری
- **WireGuard** (در حال توسعه)
- **HTTP/SOCKS** - Proxy protocols
- **Mixed** - چند پروتکله

### 📊 مدیریت و نظارت
- **Real-time Dashboard** با WebSocket
- **Traffic Statistics** دقیق
- **User Management** کامل
- **Node Management** چند سروره
- **Inbound/Outbound** مدیریت
- **Client Limits** (Traffic, IP, Expiry)
- **Prometheus Metrics**
- **Structured Logging**

### 🔒 مدیریت گواهینامه SSL
- **Self-signed** certificate generation
- **Let's Encrypt** با ACME
- **Auto-renewal** خودکار
- **Multi-domain** support
- **Certificate info** و monitoring

### 🤖 ربات تلگرام
- **Commands کامل**:
  - `/start` - شروع
  - `/status` - وضعیت سیستم
  - `/traffic <email>` - ترافیک کاربر
  - `/users` - لیست کاربران
  - `/inbounds` - لیست inbound ها
  - `/restart` - راه‌اندازی مجدد Xray
- **Daily Reports** خودکار
- **Alerts** برای رویدادهای مهم
- **Search** کاربران و کلاینت‌ها

### ⚙️ Xray Integration
- **Xray Binary** - اجرا و مدیریت خودکار Xray-core
- **gRPC Stats API** - دریافت آمار real-time از Xray
- **Config Management** - ساخت و مدیریت فایل config.json
- **Process Control** - Start/Stop/Restart Xray
- **Health Monitoring** - بررسی وضعیت و سلامت Xray

### 🔄 Subscription
- **Share Links** برای تمام پروتکل‌ها
- **QR Code** generation
- **Subscription URL** با base64
- **Auto-update** برای کلاینت‌ها
- **Multi-client** support

### 🎨 رابط کاربری
- **Modern React 18** با TypeScript
- **TailwindCSS** برای styling
- **Dark/Light Theme** با toggle
- **Responsive Design** کامل
- **Real-time Updates** با WebSocket
- **Toast Notifications**
- **Form Validation**
- **Loading States**
- **صفحات کامل شده**:
  - ✅ Dashboard - داشبورد با آمار real-time
  - ✅ Users - مدیریت کاربران با CRUD
  - ✅ Clients - مدیریت کلاینت‌ها با ترافیک
  - ✅ Inbounds - مدیریت inbound ها
  - ✅ Nodes - مدیریت سرورهای چند نود
  - ✅ Traffic - نمودار و آمار ترافیک
  - ✅ Certificates - مدیریت SSL/TLS
  - ✅ Telegram - تنظیمات ربات
  - ✅ Backup - پشتیبان‌گیری database
  - ✅ Templates - قالب‌های پیکربندی
  - ✅ Audit Logs - لاگ‌های سیستم
  - ✅ Settings - تنظیمات و پروفایل

### 🛠️ DevOps
- **Docker** multi-stage build
- **Docker Compose** آماده
- **GitHub Actions** CI/CD (simplified for v2.0.0)
- **Automated Build** verification
- **Health Checks**
- **Prometheus** monitoring
- **Grafana** dashboards (optional)

> **📝 Note**: CI pipeline در نسخه 2.0.0 به صورت ساده شده است و فقط build را تست می‌کند. تست‌های کامل در نسخه‌های بعدی اضافه خواهند شد.

### 🧪 Testing
- **Unit Tests** - تست‌های واحد برای Backend
- **Integration Tests** - تست‌های یکپارچه‌سازی
- **API Tests** - تست‌های endpoint ها
- **Frontend Tests** - React Testing Library
- **E2E Tests** - تست end-to-end (در حال توسعه)
- **Coverage Reports** - گزارش پوشش کد

### 🔒 Security Audit
- **Input Validation** - اعتبارسنجی کامل ورودی‌ها
- **SQL Injection** - محافظت با GORM/Prepared Statements
- **XSS Protection** - escape کردن output ها
- **CSRF Protection** - توکن‌های CSRF
- **Rate Limiting** - محدودسازی درخواست‌ها
- **Password Hashing** - Bcrypt با cost 12
- **JWT Security** - Secure tokens با expiry
- **HTTPS Enforcement** - اجبار استفاده از HTTPS
- **Secret Management** - مدیریت امن secrets

---

## 🔌 پروتکل‌های پشتیبانی شده

| پروتکل | وضعیت | ویژگی‌ها |
|--------|-------|----------|
| **VMess** | ✅ کامل | AlterId 0, AEAD, تمام Transports |
| **VLESS** | ✅ کامل | XTLS (Vision, Direct), Reality, Fallback |
| **Trojan** | ✅ کامل | TLS, Websocket, gRPC |
| **Shadowsocks** | ✅ کامل | تمام روش‌های رمزنگاری 2022 |
| **WireGuard** | ⚠️ در حال توسعه | - |
| **HTTP** | ✅ کامل | HTTP/SOCKS Proxy |
| **Mixed** | ✅ کامل | Multi-protocol |

### Transport Protocols
- ✅ TCP (با HTTP header)
- ✅ WebSocket (WS)
- ✅ HTTP/2 (H2)
- ✅ gRPC
- ✅ QUIC
- ✅ mKCP

### Security
- ✅ TLS 1.3
- ✅ **XTLS Vision**
- ✅ **XTLS Reality** (جدیدترین)
- ✅ Self-signed certificates
- ✅ Let's Encrypt

---

## 🚀 نصب

### پیش‌نیازها
- **Docker** 20.10+ و **Docker Compose** 2.0+
- **2GB RAM** حداقل
- **10GB** فضای دیسک
- پورت‌های **80** و **443** آزاد (برای HTTPS)

### نصب سریع (Ubuntu/Debian)

```bash
# دانلود و نصب خودکار
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

### نصب دستی

```bash
# 1. کلون پروژه
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# 2. کپی فایل تنظیمات
cp .env.example .env

# 3. تنظیم secrets (مهم!)
nano .env
# JWT_SECRET را با یک رشته 32+ کاراکتری تصادفی جایگزین کنید
# از openssl استفاده کنید: openssl rand -base64 32

# 4. اجرا با Docker Compose
docker-compose up -d

# 5. چک کردن logs
docker-compose logs -f
```

### دسترسی به پنل

پس از نصب:
- **URL**: `http://YOUR-SERVER-IP` یا `http://localhost`
- **API Docs**: `http://YOUR-SERVER-IP/docs`
- **Username**: `admin`
- **Password**: `admin`

⚠️ **خیلی مهم**: رمز عبور را بلافاصله پس از اولین ورود تغییر دهید!

---

## ⚙️ پیکربندی

### فایل `.env`

```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_MODE=production

# Database
DATABASE_TYPE=postgres
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_NAME=vui
DATABASE_USER=vui
DATABASE_PASSWORD=YOUR_SECURE_PASSWORD_HERE

# Authentication  
JWT_SECRET=YOUR_32_CHAR_RANDOM_SECRET_HERE
JWT_REFRESH_SECRET=YOUR_32_CHAR_REFRESH_SECRET_HERE
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Xray
XRAY_CONFIG_PATH=/etc/xray/config.json
XRAY_API_PORT=10085

# Telegram Bot (اختیاری)
TELEGRAM_ENABLED=true
TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN_HERE
TELEGRAM_AUTHORIZED_USERS=123456789,987654321

# Certificate (اختیاری)
CERT_DIR=/etc/v-ui/certs
CERT_ACME_EMAIL=your-email@example.com
CERT_AUTO_RENEW=true

# Monitoring (اختیاری)
ENABLE_METRICS=true
ENABLE_WEBSOCKET=true
```

### تنظیم Telegram Bot

1. با [@BotFather](https://t.me/BotFather) در تلگرام یک ربات بسازید
2. Token دریافتی را در `.env` قرار دهید
3. Chat ID خود را با [@userinfobot](https://t.me/userinfobot) بگیرید
4. Chat ID را در `TELEGRAM_AUTHORIZED_USERS` اضافه کنید
5. سرویس را restart کنید: `docker-compose restart`

---

## 📖 استفاده

### ساخت اولین Inbound

1. به پنل login کنید
2. به صفحه **Inbounds** بروید
3. روی **Create Inbound** کلیک کنید
4. تنظیمات را پر کنید:
   - **Protocol**: VLESS (توصیه می‌شود)
   - **Port**: 443
   - **Tag**: inbound-443
   - **Network**: TCP یا WS
   - **Security**: TLS یا Reality
5. **Save** کنید

### اضافه کردن Client

1. Inbound مورد نظر را باز کنید
2. **Add Client** را کلیک کنید
3. اطلاعات را وارد کنید:
   - **Email**: user@example.com
   - **Traffic Limit**: 50GB (optional)
   - **Expiry**: 30 days (optional)
   - **IP Limit**: 2 (optional)
4. **Generate** کنید

### دریافت Share Link

1. روی client کلیک کنید
2. **Get Link** را بزنید
3. لینک را copy کنید یا QR code را scan کنید
4. در کلاینت Xray خود import کنید

### Subscription

1. به صفحه **Profile** بروید
2. **Generate Subscription** را کلیک کنید
3. URL را copy کنید
4. در کلاینت خود اضافه کنید
5. کلاینت به صورت خودکار update می‌شود

---

## 🔌 API

### Authentication

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Response
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

### Inbounds

```bash
# List inbounds
curl http://localhost:8080/api/v1/inbounds \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Create inbound
curl -X POST http://localhost:8080/api/v1/inbounds \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tag": "vless-443",
    "protocol": "vless",
    "port": 443,
    "settings": { ... }
  }'
```

### Subscription

```bash
# Get subscription
curl http://localhost:8080/sub/YOUR_TOKEN

# Returns base64 encoded list of configs
```

**مستندات کامل API**: http://localhost:8080/docs

---

## 📚 مستندات

- [Installation Guide](./docs/INSTALLATION.md) - راهنمای نصب کامل
- [Configuration](./docs/CONFIGURATION.md) - تنظیمات پیشرفته
- [API Documentation](./docs/API.md) - مستندات API
- [Architecture](./docs/ARCHITECTURE.md) - معماری سیستم
- [Development](./docs/DEVELOPMENT.md) - راهنمای توسعه
- [Troubleshooting](./docs/TROUBLESHOOTING.md) - عیب‌یابی

---

## 🛠️ توسعه

### Backend (Go)

```bash
cd backend

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Run with hot reload
go run cmd/server/main.go
```

### Frontend (React)

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Build
npm run build

# Run tests
npm test
```

---

## 🤝 مشارکت

مشارکت شما خوشآمد است! لطفاً:

1. این repo را Fork کنید
2. یک branch جدید بسازید (`git checkout -b feature/amazing`)
3. تغییرات را commit کنید (`git commit -m 'Add amazing feature'`)
4. Push کنید (`git push origin feature/amazing`)
5. یک Pull Request باز کنید

---

## 📜 مجوز

این پروژه تحت مجوز **GPL-3.0** منتشر شده است. برای جزئیات بیشتر فایل [LICENSE](./LICENSE) را ببینید.

### ⚠️ هشدار قانونی

**فقط برای استفاده شخصی**. این ابزار برای اهداف قانونی privacy و امنیت طراحی شده است. هرگونه استفاده غیرقانونی اکیداً ممنوع و خلاف هدف این پروژه است. کاربران مسئول رعایت قوانین محلی خود هستند.

---

## 🙏 تشکر

- **V-UI** - پنل مدیریت مدرن Xray با معماری تمیز و امن
- الهام گرفته از: [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui)
- هسته اصلی: [Xray-core](https://github.com/XTLS/Xray-core)
- تمام مشارکت‌کنندگان open-source

---

## 📞 پشتیبانی

- **Issues**: [GitHub Issues](https://github.com/Parsa2769/V-UI2/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Parsa2769/V-UI2/discussions)
- **Telegram**: [@vui_support](https://t.me/vui_support)

---

<div align="center">

**ساخته شده با ❤️ برای جامعه ایران**

[⬆ بازگشت به بالا](#v-ui---advanced-xray-control-panel)

</div>
