# Collaboration Guide

## Tech Stack

**Backend**: Go 1.25 + Gin + GORM + PostgreSQL + Redis

**Frontend**: React 19 + Vite 8 + TypeScript + TailwindCSS 4 + shadcn/ui

## Quick Start

Prerequisites: Go 1.25+, Bun 1.0+

```bash
# Install dependencies
cd backend && go mod download
cd frontend && bun install

# Start development server
make dev

# Or start separately
make backend  # http://localhost:8000
make frontend # http://localhost:8001
```

## Common Commands

```bash
# Frontend
bun run dev          # Development server
bun run build        # Production build
bun run test         # Unit tests
bun run test:e2e     # E2E tests

# Backend
go test ./...        # Run tests

# Deploy
make deploy          # Deploy to remote server
```

## Deployment

### Initial Setup

```bash
# 1. Copy config template
cp deploy/.env.example deploy/.env

# 2. Edit config file
vim deploy/.env
```

### Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DEPLOY_HOST` | Yes | - | SSH server address |
| `DEPLOY_PORT` | No | 22 | SSH port |
| `DEPLOY_USER` | Yes | - | SSH username |
| `DEPLOY_KEY_PATH` | No | ~/.ssh/id_rsa | SSH private key path |
| `DEPLOY_REMOTE_PATH` | Yes | - | Remote deployment directory |
| `DEPLOY_BACKUP_PATH` | No | {REMOTE_PATH}-backups | Backup directory |
| `DEPLOY_KEEP_BACKUPS` | No | 5 | Number of backups to keep |
| `DEPLOY_RESTART_SERVICE` | No | true | Restart service after deploy |
| `DEPLOY_SERVICE_NAME` | No | kaleidoscope | Systemd service name |
| `DEPLOY_ENV` | No | production | Deployment environment |

### Deploy Command

```bash
make deploy
```

Deployment flow:
1. Detect remote server OS and architecture
2. Cross-compile project (backend + frontend)
3. Package as tar.gz
4. Upload and extract to server
5. Create and start systemd service

In production, the backend automatically serves frontend static files (from `./frontend` directory).

After deployment, service starts automatically. Check status:
```bash
ssh user@server.com "systemctl status kaleidoscope"
```

### SSH Key Setup

```bash
# Generate key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/kaleidoscope_deploy

# Copy public key to server
ssh-copy-id -i ~/.ssh/kaleidoscope_deploy.pub user@server.com

# Update deploy/.env
DEPLOY_KEY_PATH=~/.ssh/kaleidoscope_deploy
```

### Rollback

```bash
ssh user@server.com
cd /var/www/kaleidoscope-backups
ls -lt                                    # List backups
rm -rf /var/www/kaleidoscope/*
cp -r backup-YYYYMMDD-HHMMSS/* /var/www/kaleidoscope/
systemctl restart kaleidoscope
```

See `deploy/README.md` for details.