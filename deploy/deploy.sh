#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: $ENV_FILE not found"
    echo "Please copy .env.example to .env and configure it:"
    echo "  cp deploy/.env.example deploy/.env"
    exit 1
fi

source "$ENV_FILE"

: "${DEPLOY_HOST:?DEPLOY_HOST is required}"
: "${DEPLOY_USER:?DEPLOY_USER is required}"
: "${DEPLOY_REMOTE_PATH:?DEPLOY_REMOTE_PATH is required}"

DEPLOY_PORT="${DEPLOY_PORT:-22}"
DEPLOY_KEY_PATH="${DEPLOY_KEY_PATH:-~/.ssh/id_rsa}"
DEPLOY_BACKUP_PATH="${DEPLOY_BACKUP_PATH:-${DEPLOY_REMOTE_PATH}-backups}"
DEPLOY_KEEP_BACKUPS="${DEPLOY_KEEP_BACKUPS:-5}"
DEPLOY_RESTART_SERVICE="${DEPLOY_RESTART_SERVICE:-true}"
DEPLOY_SERVICE_NAME="${DEPLOY_SERVICE_NAME:-kaleidoscope}"
DEPLOY_ENV="${DEPLOY_ENV:-production}"
API_BASE_URL="${API_BASE_URL:-}"

SSH_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
if [ -f "$DEPLOY_KEY_PATH" ]; then
    SSH_OPTS="$SSH_OPTS -i $DEPLOY_KEY_PATH"
fi

SSH_CMD="ssh $SSH_OPTS -p $DEPLOY_PORT $DEPLOY_USER@$DEPLOY_HOST"
SCP_CMD="scp $SSH_OPTS -P $DEPLOY_PORT"

echo "=== Deployment Configuration ==="
echo "Host: $DEPLOY_HOST:$DEPLOY_PORT"
echo "User: $DEPLOY_USER"
echo "Remote Path: $DEPLOY_REMOTE_PATH"
echo "Environment: $DEPLOY_ENV"
echo "================================"

echo ""
echo "Step 1: Detecting remote OS..."
REMOTE_OS=$($SSH_CMD "uname -s" 2>/dev/null)
REMOTE_ARCH=$($SSH_CMD "uname -m" 2>/dev/null)

if [ -z "$REMOTE_OS" ]; then
    echo "Error: Cannot connect to remote server"
    exit 1
fi

echo "Remote OS: $REMOTE_OS"
echo "Remote Architecture: $REMOTE_ARCH"

case "$REMOTE_OS" in
    Linux)
        GOOS=linux
        DISTRO=$($SSH_CMD "cat /etc/os-release 2>/dev/null | grep '^ID=' | cut -d'=' -f2 | tr -d '\"'" 2>/dev/null)
        DISTRO_VERSION=$($SSH_CMD "cat /etc/os-release 2>/dev/null | grep '^VERSION_ID=' | cut -d'=' -f2 | tr -d '\"'" 2>/dev/null)
        if [ -n "$DISTRO" ]; then
            echo "Distribution: $DISTRO $DISTRO_VERSION"
        fi
        ;;
    Darwin)
        GOOS=darwin
        MACOS_VERSION=$($SSH_CMD "sw_vers -productVersion" 2>/dev/null)
        echo "macOS Version: $MACOS_VERSION"
        ;;
    *)
        echo "Error: Unsupported OS - $REMOTE_OS"
        exit 1
        ;;
esac

case "$REMOTE_ARCH" in
    x86_64|amd64)
        GOARCH=amd64
        ;;
    aarch64|arm64)
        GOARCH=arm64
        ;;
    armv7l|arm)
        GOARCH=arm
        ;;
    *)
        echo "Error: Unsupported architecture - $REMOTE_ARCH"
        exit 1
        ;;
esac

echo "Build Target: $GOOS/$GOARCH"

echo ""
echo "Step 2: Building project for $GOOS/$GOARCH..."
make build GOOS=$GOOS GOARCH=$GOARCH API_BASE_URL="$API_BASE_URL"

echo ""
echo "Step 3: Creating deployment package..."
PACKAGE_NAME="${SCRIPT_DIR}/deploy-$(date +%Y%m%d-%H%M%S).tar.gz"
tar -czf "$PACKAGE_NAME" -C build .
echo "Package created: $(basename "$PACKAGE_NAME") ($(du -h "$PACKAGE_NAME" | cut -f1))"

echo ""
echo "Step 4: Creating remote directories..."
$SSH_CMD "mkdir -p $DEPLOY_REMOTE_PATH $DEPLOY_BACKUP_PATH"

echo ""
echo "Step 5: Creating backup..."
BACKUP_NAME="backup-$(date +%Y%m%d-%H%M%S)"
$SSH_CMD "if [ -d $DEPLOY_REMOTE_PATH ] && [ \"\$(ls -A $DEPLOY_REMOTE_PATH 2>/dev/null)\" ]; then cp -r $DEPLOY_REMOTE_PATH $DEPLOY_BACKUP_PATH/$BACKUP_NAME; fi"

echo ""
echo "Step 6: Uploading package..."
$SCP_CMD "$PACKAGE_NAME" $DEPLOY_USER@$DEPLOY_HOST:/tmp/

echo ""
echo "Step 7: Extracting package..."
PACKAGE_BASENAME=$(basename "$PACKAGE_NAME")
$SSH_CMD "rm -rf $DEPLOY_REMOTE_PATH/* && tar -xzf /tmp/$PACKAGE_BASENAME -C $DEPLOY_REMOTE_PATH && rm /tmp/$PACKAGE_BASENAME"

echo ""
echo "Step 8: Cleaning old backups..."
$SSH_CMD "cd $DEPLOY_BACKUP_PATH && ls -t | tail -n +$((DEPLOY_KEEP_BACKUPS + 1)) | xargs -r rm -rf"

if [ "$DEPLOY_RESTART_SERVICE" = "true" ]; then
    echo ""
    echo "Step 9: Starting service..."
    
    $SSH_CMD "pkill -f 'kaleidoscope server' 2>/dev/null || true"
    $SSH_CMD "cd $DEPLOY_REMOTE_PATH && nohup ./kaleidoscope server > /tmp/kaleidoscope.log 2>&1 &"
    $SSH_CMD "sleep 2; pgrep -f 'kaleidoscope server' > /dev/null && echo 'Service started successfully' || echo 'Warning: Service may have failed to start'"
fi

rm -f "$PACKAGE_NAME"

echo ""
echo "✓ Deployment completed successfully!"
echo "Backup created: $BACKUP_NAME"

echo ""
echo "Checking service status..."
$SSH_CMD "ps aux | grep kaleidoscope | grep -v grep || echo 'Service status unavailable'"
