# 安装部署指南

## 环境要求

- Go 1.25+
- Bun 1.0+
- PostgreSQL 14+
- Redis 6+

## 本地开发

### 安装依赖

```bash
cd backend && go mod download
cd frontend && bun install
```

### 启动开发服务器

```bash
make dev
```

或分别启动：

```bash
make backend  # http://localhost:8000
make frontend # http://localhost:8001
```

## 生产部署

### 首次配置

1. 复制配置模板：

```bash
cp deploy/.env.example deploy/.env
```

2. 编辑 `deploy/.env`：

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `DEPLOY_HOST` | 是 | - | SSH 服务器地址 |
| `DEPLOY_PORT` | 否 | 22 | SSH 端口 |
| `DEPLOY_USER` | 是 | - | SSH 用户名 |
| `DEPLOY_KEY_PATH` | 否 | ~/.ssh/id_rsa | SSH 私钥路径 |
| `DEPLOY_REMOTE_PATH` | 是 | - | 远程部署目录 |
| `DEPLOY_BACKUP_PATH` | 否 | {REMOTE_PATH}-backups | 备份目录 |
| `DEPLOY_KEEP_BACKUPS` | 否 | 5 | 保留备份数量 |
| `DEPLOY_RESTART_SERVICE` | 否 | true | 部署后重启服务 |
| `DEPLOY_SERVICE_NAME` | 否 | kaleidoscope | Systemd 服务名 |
| `DEPLOY_ENV` | 否 | production | 部署环境 |
| `API_BASE_URL` | 否 | - | 前端 API 地址 |

### SSH 密钥配置

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/kaleidoscope_deploy
ssh-copy-id -i ~/.ssh/kaleidoscope_deploy.pub user@server.com
```

更新 `deploy/.env`：

```bash
DEPLOY_KEY_PATH=~/.ssh/kaleidoscope_deploy
```

### 部署

```bash
make deploy
```

部署流程：

1. 检测远程服务器 OS 和架构
2. 交叉编译构建项目（backend + frontend）
3. 打包成 tar.gz
4. 上传并解压到服务器
5. 创建 systemd 服务并启动

生产环境下，后端自动托管前端静态文件（`./frontend` 目录）。

查看服务状态：

```bash
ssh user@server.com "systemctl status kaleidoscope"
```

### 回滚

```bash
ssh user@server.com
cd /var/www/kaleidoscope-backups
ls -lt                                    # 查看备份列表
rm -rf /var/www/kaleidoscope/*
cp -r backup-YYYYMMDD-HHMMSS/* /var/www/kaleidoscope/
systemctl restart kaleidoscope
```

## 常用命令

### 前端

```bash
bun run dev          # 开发服务器
bun run build        # 生产构建
bun run test         # 单元测试
bun run test:e2e     # E2E 测试
```

### 后端

```bash
go test ./...        # 运行测试
make swagger         # 生成 Swagger 文档
```

### 构建

```bash
make build           # 构建到 build/ 目录
make build-backend   # 仅构建后端
make macos           # 构建 macOS 应用
```
