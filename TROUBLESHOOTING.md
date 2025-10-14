# 🔧 V-UI Troubleshooting Guide

## مشکلات رایج نصب

### 0. Docker Compose پیدا نمی‌شود (Ubuntu 25.04+)

**علائم:**
```
❌ Docker Compose plugin not found
```

**راه‌حل:**
اسکریپت جدید به صورت خودکار Docker Compose را نصب می‌کند. اگر دستی می‌خواهید نصب کنید:

```bash
# نصب دستی Docker Compose:
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose

# تست:
docker compose version
```

### 1. اسکریپت نصب گیر می‌کند

**علائم:**
- اسکریپت در "Installing system dependencies" متوقف می‌شود
- هیچ پیشرفتی دیده نمی‌شود

**راه‌حل:**
```bash
# استفاده از نسخه verbose:
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install-simple.sh)

# یا نصب دستی:
wget https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install-simple.sh
chmod +x install-simple.sh
./install-simple.sh
```

### 2. خطای Docker Permission

**علائم:**
```
permission denied while trying to connect to the Docker daemon socket
```

**راه‌حل:**
```bash
# اضافه کردن کاربر به گروه docker:
sudo usermod -aG docker $USER

# خروج و ورود مجدد یا:
newgrp docker

# تست:
docker ps
```

### 3. خطای Build Failed

**علائم:**
```
❌ Failed to build Docker images
```

**راه‌حل:**
```bash
# پاک کردن cache و build مجدد:
docker system prune -af
docker compose build --no-cache
docker compose up -d
```

### 4. پورت ۸۰ در حال استفاده است

**علائم:**
```
Bind for 0.0.0.0:80 failed: port is already allocated
```

**راه‌حل:**
```bash
# پیدا کردن سرویس استفاده‌کننده از پورت:
sudo lsof -i :80

# یا تغییر پورت در docker-compose.yml:
# ports:
#   - "8080:80"  # استفاده از پورت 8080
```

### 5. Backend سالم نیست

**علائم:**
```
Backend health check failed
```

**راه‌حل:**
```bash
# بررسی لاگ‌ها:
docker compose logs backend

# بررسی وضعیت سرویس‌ها:
docker compose ps

# ریستارت:
docker compose restart backend

# بررسی database:
docker compose logs postgres
```

### 6. Frontend نمایش داده نمی‌شود

**علائم:**
- صفحه سفید یا 404

**راه‌حل:**
```bash
# بررسی لاگ‌های frontend:
docker compose logs frontend

# بررسی nginx:
docker exec -it x3ui-frontend nginx -t

# rebuild frontend:
docker compose build frontend --no-cache
docker compose up -d frontend
```

### 7. خطای Database Connection

**علائم:**
```
Error connecting to database
```

**راه‌حل:**
```bash
# بررسی postgres:
docker compose ps postgres
docker compose logs postgres

# ریستارت database:
docker compose restart postgres

# پاک کردن و ساخت مجدد:
docker compose down -v
docker compose up -d
```

### 8. خطای npm ci

**علائم:**
```
npm ERR! code ENOLOCK
npm ERR! audit This command requires an existing lockfile
```

**راه‌حل:**
این خطا نباید اتفاق بیفتد. اگر افتاد:
```bash
cd frontend
npm install
cd ..
git add frontend/package-lock.json
git commit -m "fix: add package-lock.json"
git push
```

## دستورات مفید

### بررسی وضعیت

```bash
# وضعیت تمام سرویس‌ها:
docker compose ps

# مشاهده لاگ‌ها:
docker compose logs -f

# لاگ یک سرویس خاص:
docker compose logs -f backend

# استفاده از منابع:
docker stats
```

### پاکسازی

```bash
# توقف سرویس‌ها:
docker compose down

# پاک کردن volumes:
docker compose down -v

# پاکسازی کامل Docker:
docker system prune -af
docker volume prune -f
```

### دیباگ

```bash
# ورود به container:
docker exec -it x3ui-backend sh
docker exec -it x3ui-frontend sh

# بررسی health:
curl http://localhost:8080/health

# بررسی API:
curl http://localhost:8080/api/v1/auth/login
```

## نصب دستی

اگر اسکریپت کار نمی‌کند:

```bash
# 1. نصب Docker:
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker

# 2. Clone repository:
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# 3. تنظیم environment:
cp .env.example .env

# ویرایش .env:
nano .env
# تغییر JWT_SECRET و DATABASE_PASSWORD

# 4. اجرا:
docker compose build
docker compose up -d

# 5. بررسی:
docker compose ps
docker compose logs -f
```

## پشتیبانی

اگر مشکل حل نشد:
- 🐛 [GitHub Issues](https://github.com/Parsa2769/V-UI2/issues)
- 💬 [Discussions](https://github.com/Parsa2769/V-UI2/discussions)

لطفاً موارد زیر را در issue خود قرار دهید:
- سیستم عامل و نسخه
- نسخه Docker
- لاگ‌های خطا
- دستوراتی که اجرا کردید
