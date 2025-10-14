# 🚀 V-UI Installation Guide

## ⚡ Quick Install (One-Line Command)

### Ubuntu/Debian (Recommended)

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

This will automatically:
- ✅ Install Docker & Docker Compose
- ✅ Clone the repository
- ✅ Generate secure random secrets
- ✅ Build and start all services
- ✅ Run health checks

### After Installation

Access the panel at:
- **Web Panel**: http://YOUR-SERVER-IP
- **API Docs**: http://YOUR-SERVER-IP/docs

**Default Credentials:**
- Username: `admin`
- Password: `admin`

⚠️ **CRITICAL**: Change the default password immediately after first login!

---

## 🔧 Manual Installation

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum
- 10GB disk space
- Ports 80 and 443 available

### Steps

```bash
# 1. Clone repository
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# 2. Copy environment file
cp .env.example .env

# 3. Generate secrets
export JWT_SECRET=$(openssl rand -base64 32)
export JWT_REFRESH_SECRET=$(openssl rand -base64 32)
export DB_PASSWORD=$(openssl rand -base64 24)

# 4. Update .env file
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
sed -i "s|changeme_32_char_refresh_secret_key_here|$JWT_REFRESH_SECRET|g" .env
sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env

# 5. Build and start
docker-compose up -d

# 6. Check logs
docker-compose logs -f
```

---

## 📋 Useful Commands

```bash
# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Restart services
docker-compose restart

# Update to latest version
git pull && docker-compose up -d --build

# Backup database
docker-compose exec postgres pg_dump -U x3ui x3ui > backup.sql

# Check service status
docker-compose ps
```

---

## 🐛 Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart
```

### Can't access the panel
```bash
# Check if services are running
docker-compose ps

# Check backend health
curl http://localhost:8080/health

# Check nginx
docker-compose logs nginx
```

### Database connection issues
```bash
# Recreate database
docker-compose down -v
docker-compose up -d
```

---

## 🔒 Security Recommendations

1. **Change default password** immediately
2. **Enable 2FA** for all admin users
3. **Use strong passwords** (16+ characters)
4. **Keep secrets secure** - never share .env file
5. **Enable HTTPS** for production
6. **Regular backups** of database
7. **Update regularly** to get security patches

---

## 📚 Next Steps

After installation:
1. ✅ Change default admin password
2. ✅ Setup 2FA authentication
3. ✅ Configure SSL certificates
4. ✅ Create your first inbound
5. ✅ Add clients
6. ✅ Setup Telegram bot (optional)
7. ✅ Configure automatic backups

---

## 📞 Support

- 🐛 [Report Issues](https://github.com/Parsa2769/V-UI2/issues)
- 💬 [Discussions](https://github.com/Parsa2769/V-UI2/discussions)
- 📖 [Full Documentation](./README_V-UI.md)

---

**Made with ❤️ for the community**
