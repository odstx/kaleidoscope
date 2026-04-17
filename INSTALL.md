# Installation & Deployment Guide

## Requirements

- Go 1.25+
- Bun 1.0+
- PostgreSQL 14+
- Redis 6+

## Local Development

### Install Dependencies

```bash
cd backend && go mod download
cd frontend && bun install
```

### Start Development Server

```bash
make dev
```

Or start separately:

```bash
make backend  # http://localhost:8000
make frontend # http://localhost:8001
```

## Production Deployment

### Initial Setup

1. Copy config template:

```bash
cp deploy/.env.example deploy/.env
```

2. Edit `deploy/.env`:

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
| `API_BASE_URL` | No | - | Frontend API base URL |

### SSH Key Setup

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/kaleidoscope_deploy
ssh-copy-id -i ~/.ssh/kaleidoscope_deploy.pub user@server.com
```

Update `deploy/.env`:

```bash
DEPLOY_KEY_PATH=~/.ssh/kaleidoscope_deploy
```

### Deploy

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

Check service status:

```bash
ssh user@server.com "systemctl status kaleidoscope"
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

## Common Commands

### Frontend

```bash
bun run dev          # Development server
bun run build        # Production build
bun run test         # Unit tests
bun run test:e2e     # E2E tests
```

### Backend

```bash
go test ./...        # Run tests
make swagger         # Generate Swagger docs
```

### Build

```bash
make build           # Build to build/ directory
make build-backend   # Build backend only
make macos           # Build macOS app
```