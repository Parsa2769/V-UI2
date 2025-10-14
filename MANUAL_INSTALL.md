# 📦 V-UI Manual Installation Guide

اگر اسکریپت خودکار کار نمی‌کند، این راهنما را دنبال کنید:

## مرحله 1: نصب Docker

### Ubuntu/Debian:

```bash
# حذف نسخه‌های قدیمی:
sudo apt-get remove docker docker-engine docker.io containerd runc

# نصب پیش‌نیازها:
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# اضافه کردن GPG key:
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# اضافه کردن repository:
# برای Ubuntu 25.04 از noble استفاده کنید:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  noble stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# نصب Docker:
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# شروع Docker:
sudo systemctl daemon-reload
sudo systemctl start docker
sudo systemctl enable docker

# تست:
sudo docker run hello-world
```

## مرحله 2: نصب Docker Compose (اگر plugin نصب نشد)

```bash
# دانلود Docker Compose:
sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# اجازه اجرا:
sudo chmod +x /usr/local/bin/docker-compose

# تست:
docker-compose --version
# یا
docker compose version
```

## مرحله 3: کلون پروژه

```bash
# نصب git:
sudo apt-get install -y git

# کلون repository:
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2
```

## مرحله 4: تنظیم Environment

```bash
# کپی فایل نمونه:
cp .env.example .env

# ویرایش فایل:
nano .env

# یا تولید خودکار secrets:
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
JWT_REFRESH_SECRET=$(openssl rand -base64 32 | tr -d '\n')
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')

sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
sed -i "s|changeme_32_char_refresh_secret_key_here|$JWT_REFRESH_SECRET|g" .env
sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env
```

## مرحله 5: Build و اجرا

```bash
# Build images (5-10 دقیقه طول می‌کشد):
docker compose build

# یا با docker-compose:
docker-compose build

# شروع سرویس‌ها:
docker compose up -d

# یا:
docker-compose up -d

# بررسی وضعیت:
docker compose ps
docker compose logs -f
```

## مرحله 6: دسترسی به پنل

بعد از نصب موفق:

- **پنل وب**: http://YOUR-SERVER-IP
- **API Docs**: http://YOUR-SERVER-IP/docs
- **Username**: admin
- **Password**: admin

⚠️ **حتماً رمز عبور را تغییر دهید!**

## عیب‌یابی

### Docker daemon شروع نمی‌شود:

```bash
# بررسی وضعیت:
sudo systemctl status docker

# بررسی لاگ‌ها:
sudo journalctl -xeu docker.service

# ریستارت:
sudo systemctl restart docker
```

### Build خطا می‌دهد:

```bash
# پاک کردن cache:
docker system prune -af
docker volume prune -f

# Build مجدد:
docker compose build --no-cache
```

### سرویس‌ها شروع نمی‌شوند:

```bash
# بررسی لاگ‌ها:
docker compose logs backend
docker compose logs frontend
docker compose logs postgres

# بررسی پورت‌ها:
sudo lsof -i :80
sudo lsof -i :8080

# ریستارت:
docker compose down
docker compose up -d
```

### Permission Denied:

```bash
# اضافه کردن user به گروه docker:
sudo usermod -aG docker $USER

# خروج و ورود مجدد یا:
newgrp docker

# تست:
docker ps
```

## دستورات مفید

```bash
# مشاهده لاگ‌ها:
docker compose logs -f

# توقف:
docker compose down

# ریستارت:
docker compose restart

# ورود به container:
docker exec -it x3ui-backend sh
docker exec -it x3ui-frontend sh

# بررسی استفاده از منابع:
docker stats

# پاک کردن کامل:
docker compose down -v
docker system prune -af
```

## نصب از Source (توسعه‌دهندگان)

### Backend:

```bash
cd backend

# نصب Go 1.21:
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# نصب dependencies:
go mod download

# اجرا:
go run cmd/server/main.go
```

### Frontend:

```bash
cd frontend

# نصب Node.js 18:
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# نصب dependencies:
npm install

# اجرا در حالت development:
npm run dev

# Build برای production:
npm run build
```

## پشتیبانی

- 🐛 [GitHub Issues](https://github.com/Parsa2769/V-UI2/issues)
- 💬 [Discussions](https://github.com/Parsa2769/V-UI2/discussions)
- 📖 [Documentation](./README_V-UI.md)
