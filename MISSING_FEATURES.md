# ویژگی‌های باقی‌مانده برای کامل شدن

## 🎯 اولویت بالا (Critical)

### 1. Xray Core Integration
**Status**: نیمه‌کاره (فقط ساختار اولیه)

باید پیاده‌سازی بشه:
```go
// backend/internal/xray/manager.go
- InstallXray() - نصب Xray core
- UpdateXray() - بروزرسانی Xray
- RestartXray() - Restart service
- GetXrayStats() - دریافت آمار از Xray API
- ApplyConfig() - اعمال config جدید
```

**فایل‌های مورد نیاز:**
- `backend/internal/xray/manager.go`
- `backend/internal/xray/stats.go`
- `backend/internal/xray/config_builder.go`

### 2. Protocol Handlers
**Status**: ندارم

باید برای هر پروتکل handler بسازیم:

```go
// backend/internal/xray/protocols/
- vmess.go      // VMESS config generator
- vless.go      // VLESS config generator
- trojan.go     // Trojan config generator
- shadowsocks.go // Shadowsocks config generator
- wireguard.go  // WireGuard (اگر نیاز باشه)
```

**هر handler باید:**
- Config generation
- Validation
- QR code data
- Share link generation

### 3. Inbound/Outbound Management
**Status**: Model داره، handler نداره

```go
// backend/internal/service/inbound.go
type InboundService struct {
    CreateInbound(protocol, port, settings)
    UpdateInbound(id, settings)
    DeleteInbound(id)
    ListInbounds()
    GetInboundStats(id)
}

// backend/internal/api/inbound_handler.go
- POST /api/v1/inbounds
- GET /api/v1/inbounds
- GET /api/v1/inbounds/:id
- PUT /api/v1/inbounds/:id
- DELETE /api/v1/inbounds/:id
```

**Models مورد نیاز:**
```go
type Inbound struct {
    ID       uuid.UUID
    Tag      string
    Protocol string // vmess, vless, trojan, etc.
    Port     int
    Settings json.RawMessage
    Enable   bool
}

type Outbound struct {
    ID       uuid.UUID
    Tag      string
    Protocol string
    Settings json.RawMessage
}
```

### 4. QR Code & Subscription
**Status**: نداریم

```go
// backend/internal/service/subscription.go
- GenerateQRCode(clientID) -> QR image
- GenerateShareLink(clientID) -> vmess:// or vless://
- GenerateSubscriptionLink(userID) -> subscription URL
- GetSubscriptionContent(token) -> base64 encoded configs
```

**Endpoints:**
- `GET /api/v1/clients/:id/qr`
- `GET /api/v1/clients/:id/link`
- `GET /sub/:token` (subscription endpoint)

### 5. Real Traffic Tracking
**Status**: ساختار داره، integration نداره

باید با Xray Stats API ارتباط بگیریم:

```go
// backend/internal/xray/stats.go
type StatsClient struct {
    conn *grpc.ClientConn
}

func (s *StatsClient) GetUserTraffic(email string) (up, down int64)
func (s *StatsClient) GetInboundTraffic(tag string) (up, down int64)
func (s *StatsClient) ResetUserTraffic(email string)
```

**Background Job:**
```go
// هر 10 ثانیه یکبار
- دریافت traffic از Xray
- آپدیت database
- چک کردن limits
- قطع کردن کاربرانی که limit رو رد کردن
```

---

## 🎯 اولویت متوسط (Important)

### 6. Certificate Management
**Status**: نداریم

```go
// backend/internal/cert/manager.go
- GenerateSelfSigned()
- RequestLetsEncrypt(domain)
- RenewCertificate()
- GetCertificateInfo()
```

**Frontend:**
- صفحه SSL/TLS Settings
- Auto-renewal toggle
- Certificate upload

### 7. Routing Rules
**Status**: نداریم

```go
// backend/internal/xray/routing.go
type RoutingRule struct {
    Type        string   // field
    Domain      []string
    IP          []string
    Port        string
    Network     string
    Protocol    []string
    OutboundTag string
}

func (r *RoutingService) AddRule(rule RoutingRule)
func (r *RoutingService) UpdateRule(id, rule)
func (r *RoutingService) DeleteRule(id)
```

**Frontend:**
- صفحه Routing Rules
- Domain blocking/allowing
- IP whitelist/blacklist
- Protocol routing

### 8. Notifications
**Status**: نداریم

```go
// backend/internal/notification/
- telegram.go  // Telegram bot
- email.go     // Email sender
- webhook.go   // Webhook calls

type NotificationService struct {
    SendAlert(title, message string)
    SendTrafficWarning(client)
    SendExpiryWarning(client)
    SendSystemAlert(error)
}
```

**تنظیمات در .env:**
```env
TELEGRAM_BOT_TOKEN=
TELEGRAM_CHAT_ID=
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
```

### 9. Backup & Restore
**Status**: فقط database backup داریم

باید اضافه بشه:
- Xray config backup
- Full system backup
- One-click restore
- Scheduled backups
- Cloud backup (S3, etc.)

```go
// backend/internal/backup/manager.go
func BackupAll() (file string, error)
func RestoreFromBackup(file string) error
func ScheduledBackup(interval time.Duration)
```

### 10. Migration Tool
**Status**: نداریم

```go
// cmd/migrate/main.go
- ImportFromOldXUI(dbPath string)
- ImportUsers()
- ImportInbounds()
- ImportClients()
- ConvertOldConfig()
```

---

## 🎯 اولویت پایین (Nice to Have)

### 11. Telegram Bot
```go
// backend/internal/bot/telegram.go
Commands:
- /start
- /status
- /traffic
- /users
- /restart
```

### 12. Multi-language
- فارسی ✅ (نیمه)
- English ✅ (نیمه)
- عربی ❌
- 中文 ❌
- Русский ❌

### 13. Advanced Dashboard
- Traffic charts (daily/weekly/monthly)
- System resource monitoring
- Geographic traffic map
- Protocol distribution pie chart

### 14. User Portal
صفحه جداگانه برای کاربران:
- مشاهده traffic خودشون
- دانلود config
- QR code
- تاریخ انقضا

### 15. API Rate Limiting per User
- هر کاربر limit مجزا
- Custom rate limits per role

---

## 📦 Packages مورد نیاز برای تکمیل

### Go Packages
```go
// go.mod
require (
    // QR Code
    github.com/skip2/go-qrcode v0.0.0
    
    // Telegram Bot
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    
    // Email
    github.com/jordan-wright/email v4.0.1
    
    // gRPC برای Xray Stats
    google.golang.org/grpc v1.60.0
    
    // SSL/TLS
    golang.org/x/crypto/acme/autocert v0.0.0
    
    // Image processing
    github.com/disintegration/imaging v1.6.2
)
```

### Frontend Packages
```json
{
  "dependencies": {
    "qrcode.react": "^3.1.0",
    "react-chartjs-2": "^5.2.0",
    "chart.js": "^4.4.0",
    "leaflet": "^1.9.4",
    "react-leaflet": "^4.2.1"
  }
}
```

---

## 🗂️ ساختار فایل‌های باقی‌مانده

```
backend/
├── internal/
│   ├── xray/
│   │   ├── manager.go        ❌ نداریم
│   │   ├── stats.go          ❌ نداریم
│   │   ├── config_builder.go ❌ نداریم
│   │   ├── protocols/
│   │   │   ├── vmess.go      ❌ نداریم
│   │   │   ├── vless.go      ❌ نداریم
│   │   │   ├── trojan.go     ❌ نداریم
│   │   │   └── ss.go         ❌ نداریم
│   │   └── routing.go        ❌ نداریم
│   ├── service/
│   │   ├── inbound.go        ❌ نداریم
│   │   ├── outbound.go       ❌ نداریم
│   │   ├── subscription.go   ❌ نداریم
│   │   └── backup.go         ❌ نداریم (فقط script داریم)
│   ├── notification/
│   │   ├── telegram.go       ❌ نداریم
│   │   ├── email.go          ❌ نداریم
│   │   └── webhook.go        ❌ نداریم
│   ├── cert/
│   │   └── manager.go        ❌ نداریم
│   └── bot/
│       └── telegram.go       ❌ نداریم
├── cmd/
│   └── migrate/
│       └── main.go           ❌ نداریم

frontend/
├── src/
│   ├── pages/
│   │   ├── InboundsPage.tsx  ❌ نداریم
│   │   ├── RoutingPage.tsx   ❌ نداریم
│   │   ├── CertsPage.tsx     ❌ نداریم
│   │   └── UserPortal.tsx    ❌ نداریم
│   └── components/
│       ├── QRCode.tsx        ❌ نداریم
│       ├── TrafficChart.tsx  ❌ نداریم
│       └── GeoMap.tsx        ❌ نداریم
```

---

## 📝 مراحل تکمیل (به ترتیب اولویت)

### Phase 1: Core Xray (2-3 روز)
1. ✅ Xray manager implementation
2. ✅ Protocol handlers (VMESS, VLESS, Trojan)
3. ✅ Stats API integration
4. ✅ Traffic tracking background job

### Phase 2: Client Features (1-2 روز)
5. ✅ QR code generation
6. ✅ Share link generation
7. ✅ Subscription system
8. ✅ Inbound/Outbound management

### Phase 3: Advanced Features (2-3 روز)
9. ✅ Routing rules
10. ✅ Certificate management
11. ✅ Notifications (Telegram, Email)
12. ✅ Enhanced backup/restore

### Phase 4: Polish (1-2 روز)
13. ✅ Migration tool
14. ✅ User portal
15. ✅ Advanced charts
16. ✅ Testing

**تخمین زمان کل**: 6-10 روز کار focused

---

## 🎓 منابع برای پیاده‌سازی

1. **Xray API Docs**: https://xtls.github.io/
2. **Original 3x-ui Code**: https://github.com/MHSanaei/3x-ui
3. **Xray Config Examples**: https://github.com/XTLS/Xray-examples
4. **gRPC Stats API**: Xray stats service protocol

---

## ✅ خلاصه

**آنچه داریم (60%)**:
- ✅ Infrastructure کامل
- ✅ Backend architecture
- ✅ Frontend مدرن
- ✅ Authentication & Authorization
- ✅ Database & Models
- ✅ Basic CRUD operations
- ✅ DevOps & CI/CD

**آنچه نداریم (40%)**:
- ❌ Xray core integration کامل
- ❌ Protocol handlers
- ❌ QR codes & subscriptions
- ❌ Real traffic tracking
- ❌ Routing rules
- ❌ Certificates
- ❌ Notifications
- ❌ Advanced features

**نتیجه**: پروژه یک **foundation عالی** هست اما برای استفاده واقعی نیاز به 6-10 روز کار اضافه داره.
