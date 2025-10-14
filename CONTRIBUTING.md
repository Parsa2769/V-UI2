# Contributing to 3X-UI Modern

Thank you for considering contributing to 3X-UI Modern!

## How to Contribute

### Reporting Bugs
- Check existing issues first
- Use the bug report template
- Include OS, version, and reproduction steps

### Pull Requests

**Process:**
1. Fork repo and create branch from `main`
2. Add tests for new code
3. Update documentation
4. Ensure tests pass
5. Submit PR

**PR Checklist:**
- [ ] Code follows project style
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] CI passes

### Development Setup

```bash
# Clone and setup
git clone https://github.com/YOUR_USERNAME/3x-ui-modern.git
cd 3x-ui-modern

# Backend
cd backend && go mod download

# Frontend
cd frontend && npm install

# Run tests
go test ./...
npm test
```

### Commit Messages

Format: `type(scope): description`

Types: feat, fix, docs, style, refactor, test, chore

Examples:
- `feat(ui): add dark mode toggle`
- `fix(auth): prevent token reuse`
- `docs(api): update swagger specs`

### Code Style

**Go:**
- Run `golangci-lint run`
- Follow standard Go conventions
- Add comments for exported functions

**TypeScript:**
- Run `npm run lint`
- Use functional components
- Type everything explicitly

## Questions?

Open a discussion or contact maintainers.
