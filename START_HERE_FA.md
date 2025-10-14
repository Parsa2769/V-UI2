# 🚀 از اینجا شروع کنید - START HERE

## ✅ پروژه کامل و آماده برای استفاده است!

این یک fork مدرن و کامل از پروژه 3x-ui هست که با تکنولوژی‌های جدید و امکانات پیشرفته ساخته شده.

---

## 📋 قبل از شروع

### چک کنید همه چیز کامل هست:

```bash
# چک کردن ساختار پروژه
ls -la

# باید این فایل‌ها رو ببینید:
✅ README.md (مستندات کامل)
✅ docker-compose.yml (برای اجرا)
✅ .env.example (تنظیمات)
✅ install-ubuntu.sh (نصب اتوماتیک)
✅ backend/ (کد Go)
✅ frontend/ (کد React)
```

---

## 🎯 دو راه برای استفاده

### 1️⃣ اجرای محلی (برای تست و توسعه)

```bash
# کپی کردن فایل تنظیمات
copy .env.example .env

# ویرایش تنظیمات (در ویندوز)
notepad .env

# اجرا با Docker
docker-compose up -d

# مشاهده logs
docker-compose logs -f

# دسترسی
http://localhost:8080
```

**لاگین پیش‌فرض:**
- Username: `admin`
- Password: `admin`

---

### 2️⃣ نصب روی سرور Ubuntu (پروداکشن)

#### الف) نصب اتوماتیک (توصیه می‌شود)

```bash
# روی سرور Ubuntu اجرا کنید:
bash <(curl -fsSL https://raw.githubusercontent.com/YOUR-USERNAME/3x-ui-modern/main/install-ubuntu.sh)
```

#### ب) نصب دستی

```bash
# 1. کلون پروژه
git clone https://github.com/YOUR-USERNAME/3x-ui-modern.git
cd 3x-ui-modern

# 2. اجرای اسکریپت نصب
chmod +x install-ubuntu.sh
./install-ubuntu.sh

# 3. صبر کنید تا نصب کامل بشه (5-10 دقیقه)

# 4. دسترسی به پنل
http://YOUR-SERVER-IP
```

---

## 📤 آپلود به GitHub

### مرحله 1: ساخت Repository

1. برو به: https://github.com/new
2. Repository name: `3x-ui-modern`
3. Visibility: Public یا Private
4. **Initialize with README را تیک نزن**
5. Create repository

### مرحله 2: آپلود کد

```bash
# در پوشه پروژه:
cd c:\Users\parsa\CascadeProjects\windsurf-project-2

# Initialize git
git init

# Add همه فایل‌ها
git add .

# اولین commit
git commit -m "Initial commit: 3X-UI Modern v2.0.0"

# اتصال به GitHub (URL خودتون رو جایگزین کنید)
git remote add origin https://github.com/YOUR-USERNAME/3x-ui-modern.git

# Push
git branch -M main
git push -u origin main
```

### مرحله 3: بروزرسانی URL ها

بعد از آپلود، این فایل‌ها رو update کنید:
- `README.md` → جایگزین `yourusername` با username واقعی GitHub
- `install-ubuntu.sh` → جایگزین URL

```bash
git add README.md install-ubuntu.sh
git commit -m "docs: update GitHub URLs"
git push
```

راهنمای کامل: `GITHUB_SETUP.md`

---

## 🛠️ تنظیمات مهم

### فایل `.env`

**این تنظیمات رو حتماً تغییر بدید:**

```env
# امنیتی (خیلی مهم!)
JWT_SECRET=تولید_یک_رشته_32_کاراکتری_تصادفی
JWT_REFRESH_SECRET=تولید_یک_رشته_32_کاراکتری_تصادفی_دیگر
DATABASE_PASSWORD=یک_پسورد_قوی

# سرور
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# دیتابیس
DATABASE_TYPE=postgres
DATABASE_HOST=postgres
DATABASE_NAME=x3ui
```

**تولید secret های امن:**
```bash
# در Linux/Mac:
openssl rand -base64 32

# در PowerShell:
-join ((65..90) + (97..122) + (48..57) | Get-Random -Count 32 | ForEach-Object {[char]$_})
```

---

## 📊 بررسی وضعیت نصب

### چک کردن سرویس‌ها

```bash
# وضعیت containers
docker-compose ps

# باید ببینید:
✅ postgres - healthy
✅ backend - healthy
✅ nginx - running
```

### چک کردن سلامت سیستم

```bash
# Health check
curl http://localhost:8080/health

# Response باید باشه:
{"status":"healthy","time":1234567890}

# API test
curl http://localhost:8080/api/v1/metrics/summary

# Metrics (Prometheus)
curl http://localhost:8080/metrics
```

---

## 🎨 ویژگی‌های پروژه

### Backend (Go)
✅ JWT Authentication  
✅ 2FA (TOTP)  
✅ RBAC (3 نقش)  
✅ RESTful API  
✅ WebSocket  
✅ PostgreSQL/SQLite  
✅ Prometheus Metrics  
✅ Structured Logging  
✅ Xray Integration  

### Frontend (React)
✅ TypeScript  
✅ TailwindCSS  
✅ Dark/Light Theme  
✅ Responsive Design  
✅ Real-time Updates  
✅ Form Validation  
✅ Toast Notifications  

### DevOps
✅ Docker Compose  
✅ Multi-stage Build  
✅ GitHub Actions CI/CD  
✅ Automated Tests  
✅ Backup Scripts  
✅ Health Checks  

---

## 📚 مستندات

- **README.md** → مستندات کامل (فارسی + انگلیسی)
- **QUICKSTART.md** → راهنمای سریع
- **FINAL_CHECKLIST.md** → چک‌لیست نهایی
- **GITHUB_SETUP.md** → راهنمای آپلود به GitHub
- **docs/API.md** → مستندات API
- **docs/DEPLOYMENT.md** → راهنمای استقرار
- **docs/ARCHITECTURE.md** → معماری سیستم

---

## 🔒 امنیت

### بعد از نصب حتماً:
1. ✅ پسورد admin رو تغییر بدید
2. ✅ JWT secrets رو تغییر بدید
3. ✅ 2FA رو فعال کنید
4. ✅ HTTPS تنظیم کنید
5. ✅ Firewall کانفیگ کنید
6. ✅ پشتیبان‌گیری منظم

---

## 🆘 مشکلات رایج

### سرویس‌ها start نمی‌شن
```bash
docker-compose down
docker-compose up -d --force-recreate
docker-compose logs -f
```

### پورت اشغال است
```bash
# در .env تغییر بدید:
SERVER_PORT=8081
```

### دسترسی به Docker نداره
```bash
# User رو به گروه docker اضافه کنید
sudo usermod -aG docker $USER
# logout و login مجدد
```

### Database connection error
```bash
# Reset کردن volumes
docker-compose down -v
docker-compose up -d
```

---

## 📞 پشتیبانی

- **GitHub Issues**: برای گزارش bug
- **GitHub Discussions**: برای سوالات
- **Documentation**: مستندات کامل در `docs/`

---

## ⚖️ قانونی

**فقط برای استفاده شخصی**. این ابزار برای اهداف قانونی طراحی شده. استفاده غیرقانونی اکیداً ممنوع است.

License: GPL-3.0

---

## ✨ نسخه

**Version**: 2.0.0  
**Status**: ✅ Production Ready  
**Last Updated**: 2024-10-14

---

## 🎉 آماده هستید!

همه چیز آماده است. می‌تونید:
1. روی سرور نصب کنید
2. به GitHub آپلود کنید
3. استفاده کنید

**موفق باشید! 🚀**

---

## دستورات سریع

```bash
# اجرا
docker-compose up -d

# Logs
docker-compose logs -f

# Restart
docker-compose restart

# Stop
docker-compose down

# Backup
./scripts/backup.sh

# Update
git pull && docker-compose up -d --build
```

---

**💡 نکته**: اگر اولین باره که با Docker کار می‌کنید، حتماً مستندات Docker رو بخونید.
