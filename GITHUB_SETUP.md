# راهنمای آپلود پروژه به GitHub

## مرحله 1: ساخت Repository در GitHub

1. به https://github.com بروید و login کنید
2. روی + در گوشه بالا سمت راست کلیک کنید
3. "New repository" را انتخاب کنید
4. تنظیمات:
   - **Repository name**: `V-UI2`
   - **Description**: `Modern Xray Management Panel - Fork of 3x-ui with enhanced UI and features`
   - **Visibility**: Public یا Private (به دلخواه)
   - **❌ Initialize with README را تیک نزنید** (چون ما README داریم)
5. "Create repository" را کلیک کنید

## مرحله 2: آماده‌سازی فایل‌ها

در پوشه پروژه این دستورات رو اجرا کنید:

```bash
cd c:\Users\parsa\CascadeProjects\windsurf-project-2

# مطمئن شوید .env در .gitignore هست (هست!)
# مطمئن شوید data/ و backups/ در .gitignore هست (هست!)
```

## مرحله 3: Initialize Git Repository

```bash
# Git repository رو initialize کنید
git init

# همه فایل‌ها رو add کنید
git add .

# اولین commit
git commit -m "Initial commit: V-UI v2.0.0

- Complete backend with Go 1.21
- Modern React 18 frontend with TypeScript
- JWT authentication + 2FA + RBAC
- All 13 pages fully implemented
- Docker Compose ready
- Comprehensive documentation
- CI/CD pipelines
- Production ready"
```

## مرحله 4: Connect به GitHub Repository

Repository URL خودتون رو جایگزین کنید:

```bash
# Remote repository رو اضافه کنید
git remote add origin https://github.com/Parsa2769/V-UI2.git

# Branch اصلی رو به main تغییر بدید
git branch -M main

# Push کنید
git push -u origin main
```

## مرحله 5: تنظیمات Repository در GitHub

بعد از push کردن، به صفحه repository در GitHub برید و:

### 1. About بخش
- **Description**: Modern Xray Management Panel - Fork of 3x-ui with enhanced UI, JWT auth, 2FA, RBAC, and production-ready features
- **Website**: https://github.com/Parsa2769/V-UI2
- **Topics** اضافه کنید:
  - `xray`
  - `v2ray`
  - `vpn`
  - `proxy`
  - `panel`
  - `go`
  - `react`
  - `typescript`
  - `docker`
  - `postgresql`

### 2. Branches Protection (اختیاری ولی توصیه میشه)
Settings → Branches → Add rule:
- **Branch name pattern**: `main`
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass before merging
- Save changes

### 3. GitHub Actions
- به tab "Actions" برید
- Workflows اتوماتیک فعال میشن
- CI pipeline بعد از هر push اجرا میشه

### 4. Releases
برای ساخت اولین release:
1. به tab "Releases" برید
2. "Create a new release" کلیک کنید
3. Tag version: `v2.0.0`
4. Release title: `V-UI v2.0.0 - Initial Release`
5. Description رو از CHANGELOG.md کپی کنید
6. "Publish release" کلیک کنید

## مرحله 6: بروزرسانی README.md

فایل README.md رو باز کنید و این موارد رو update کنید:

```markdown
# تمام لینک‌های yourusername به Parsa2769/V-UI2 تغییر داده شد

# دستور نصب تک‌خطی:
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

## مرحله 7: Commit تغییرات README

```bash
git add README.md install-ubuntu.sh
git commit -m "docs: update GitHub URLs"
git push
```

## مرحله 8: تست Installation Script

روی یک سرور تمیز Ubuntu تست کنید:

```bash
# روی سرور Ubuntu
bash <(curl -Ls https://raw.githubusercontent.com/Parsa2769/V-UI2/main/install.sh)
```

## دستورات مفید برای بعداً

### اضافه کردن تغییرات جدید
```bash
git add .
git commit -m "feat: add new feature"
git push
```

### ساخت branch جدید
```bash
git checkout -b feature/new-feature
# کارهاتون رو انجام بدید
git add .
git commit -m "feat: new feature description"
git push -u origin feature/new-feature
```

### Merge کردن با main
```bash
git checkout main
git merge feature/new-feature
git push
```

### Tag کردن version جدید
```bash
git tag -a v2.0.1 -m "Version 2.0.1"
git push origin v2.0.1
```

## چک‌لیست قبل از Push

- [ ] `.env` در .gitignore هست
- [ ] `data/` و `backups/` در .gitignore هست  
- [ ] SSL certificates در .gitignore هست
- [ ] README.md update شده با URL های درست
- [ ] CHANGELOG.md کامل هست
- [ ] License فایل وجود داره
- [ ] همه tests pass میشن
- [ ] Documentation کامل هست

## نکات امنیتی مهم

⚠️ **هیچ وقت این موارد رو push نکنید:**
- فایل `.env` با secrets واقعی
- Database backups
- SSL certificates
- API keys
- Passwords
- Private keys
- Production data

## مشکلات رایج

### Permission denied (publickey)
```bash
# SSH key اضافه کنید:
ssh-keygen -t ed25519 -C "your_email@example.com"
# Public key رو به GitHub اضافه کنید
```

### Large files
```bash
# از Git LFS استفاده کنید:
git lfs install
git lfs track "*.zip"
git add .gitattributes
```

### Remove sensitive file from history
```bash
# اگر اشتباهی فایل حساس رو push کردید:
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch PATH-TO-FILE" \
  --prune-empty --tag-name-filter cat -- --all
git push origin --force --all
```

## بعد از Setup

1. **README Badge ها** رو اضافه کنید:
   - CI status
   - License
   - Version
   - Stars

2. **GitHub Sponsors** فعال کنید (اختیاری)

3. **Security Policy** اضافه کنید:
   - SECURITY.md فایل بسازید

4. **Code of Conduct** اضافه کنید

5. **Wiki** رو setup کنید با documentation بیشتر

## Resources

- GitHub Docs: https://docs.github.com
- Git Docs: https://git-scm.com/doc
- GitHub Actions: https://docs.github.com/en/actions

---

**موفق باشید! 🚀**

اگر مشکلی پیش اومد، به مستندات GitHub یا Stack Overflow مراجعه کنید.
