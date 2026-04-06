# 协作文档

## 技术栈

**后端**: Go 1.25 + Gin + GORM + PostgreSQL + Redis

**前端**: React 19 + Vite 8 + TypeScript + TailwindCSS 4 + shadcn/ui

## 快速开始

前置要求：Go 1.25+, Bun 1.0+

```bash
# 安装依赖
cd backend && go mod download
cd frontend && bun install

# 启动开发服务器
make dev

# 或分别启动
make backend  # http://localhost:8000
make frontend # http://localhost:8001
```

## 常用命令

```bash
# 前端
bun run dev          # 开发服务器
bun run build        # 生产构建
bun run test         # 单元测试
bun run test:e2e     # E2E 测试

# 后端
go test ./...        # 运行测试

# 部署
make deploy          # 部署到远程服务器
```

## 部署

### 首次配置

```bash
# 1. 复制配置模板
cp deploy/.env.example deploy/.env

# 2. 编辑配置文件
vim deploy/.env
```

### 配置说明

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

### 部署命令

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

部署完成后，服务自动启动。查看状态：
```bash
ssh user@server.com "systemctl status kaleidoscope"
```

### SSH 密钥配置

```bash
# 生成密钥
ssh-keygen -t rsa -b 4096 -f ~/.ssh/kaleidoscope_deploy

# 复制公钥到服务器
ssh-copy-id -i ~/.ssh/kaleidoscope_deploy.pub user@server.com

# 更新 deploy/.env
DEPLOY_KEY_PATH=~/.ssh/kaleidoscope_deploy
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

详细文档见 `deploy/README.md`
