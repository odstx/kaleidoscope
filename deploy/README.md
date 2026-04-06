# Deployment Guide

## Quick Start

1. Copy the example configuration:
   ```bash
   cp deploy/.env.example deploy/.env
   ```

2. Edit `deploy/.env` with your server details:
   ```bash
   vim deploy/.env
   ```

3. Deploy:
   ```bash
   make deploy
   ```

## Configuration

Edit `deploy/.env` with the following variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DEPLOY_HOST` | Yes | - | SSH server hostname or IP |
| `DEPLOY_PORT` | No | 22 | SSH port |
| `DEPLOY_USER` | Yes | - | SSH username |
| `DEPLOY_KEY_PATH` | No | ~/.ssh/id_rsa | Path to SSH private key |
| `DEPLOY_REMOTE_PATH` | Yes | - | Remote deployment directory |
| `DEPLOY_BACKUP_PATH` | No | {REMOTE_PATH}-backups | Backup directory |
| `DEPLOY_KEEP_BACKUPS` | No | 5 | Number of backups to keep |
| `DEPLOY_RESTART_SERVICE` | No | true | Restart service after deploy |
| `DEPLOY_SERVICE_NAME` | No | kaleidoscope | Systemd service name |
| `DEPLOY_ENV` | No | production | Deployment environment |

## SSH Key Setup

1. Generate an SSH key (if you don't have one):
   ```bash
   ssh-keygen -t rsa -b 4096 -f ~/.ssh/kaleidoscope_deploy
   ```

2. Copy the public key to the server:
   ```bash
   ssh-copy-id -i ~/.ssh/kaleidoscope_deploy.pub user@server.com
   ```

3. Update `deploy/.env`:
   ```bash
   DEPLOY_KEY_PATH=~/.ssh/kaleidoscope_deploy
   ```

## Deployment Process

The deployment script performs these steps:

1. **Detect**: Detects remote server OS and architecture
2. **Build**: Cross-compiles backend and frontend for target platform
3. **Package**: Creates a tar.gz archive of the build (includes config.yaml)
4. **Upload**: Transfers package via SCP and extracts on server
5. **Backup**: Creates a timestamped backup of previous deployment
6. **Cleanup**: Removes old backups (keeps last N backups)
7. **Service**: Creates/updates systemd service and starts it

## Deployed Files

The following files are deployed to the server:

- `kaleidoscope` - Backend executable
- `frontend/` - Frontend static files
- `config.yaml` - Application configuration (from `backend/config/config.yaml`)

To customize production config, edit `backend/config/config.yaml` before deploying.

After deployment, the service runs automatically. Check status:
```bash
ssh user@server.com "systemctl status kaleidoscope"
```

## Rollback

To rollback to a previous version:

```bash
ssh user@server.com
cd /var/www/kaleidoscope-backups
ls -lt  # List available backups
rm -rf /var/www/kaleidoscope/*
cp -r backup-YYYYMMDD-HHMMSS/* /var/www/kaleidoscope/
systemctl restart kaleidoscope
```

## Troubleshooting

### Permission Denied

Ensure:
- SSH key has correct permissions: `chmod 600 ~/.ssh/id_rsa`
- Remote user has write permissions to `DEPLOY_REMOTE_PATH`
- Remote user can restart the service (sudo or service ownership)

### Connection Timeout

Check:
- Server is accessible: `ping your-server.com`
- SSH port is open: `nc -zv your-server.com 22`
- Firewall allows SSH connections

### Service Won't Start

Check logs on the server:
```bash
ssh user@server.com 'journalctl -u kaleidoscope -n 50'
```
