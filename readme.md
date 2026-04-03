# Kaleidoscope

全栈应用，Go 后端 + React 前端，分层架构设计。

## 架构

```
┌─────────────────────────────────────────────────────────┐
│                    对客层 (Port 8001)                    │
│              Web Frontend (React + TailwindCSS)          │
└─────────────────────────────────────────────────────────┘
                            ↓ HTTP/JSON
┌─────────────────────────────────────────────────────────┐
│                    接口层 (Port 8000)                    │
│                   HTTP API Server (Gin)                 │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                      应用层                              │
│   Controllers → Services → Models (业务逻辑处理)        │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                      数据层                              │
│            PostgreSQL + Redis + GORM                    │
└─────────────────────────────────────────────────────────┘
```

## 技术栈

### 后端
- Go 1.25
- Gin - HTTP 框架
- GORM - ORM
- PostgreSQL - 数据库
- Redis - 缓存
- Viper - 配置管理
- Zap - 日志
- Cobra - CLI
- Swagger - API 文档

### 前端
- React 19
- Vite 8
- TypeScript
- TailwindCSS 4 + shadcn/ui
- React Router v7
- react-hook-form + zod
- i18next - 国际化
- Vitest + Playwright - 测试

## 项目结构

```
.
├── backend/                 # Go 后端
│   ├── cmd/                # CLI 命令
│   ├── config/             # 配置
│   ├── controllers/        # 控制器
│   ├── database/           # 数据库连接
│   ├── docs/               # Swagger 文档
│   ├── middleware/         # 中间件
│   ├── models/             # 模型
│   ├── server/             # 服务器配置
│   ├── services/           # 服务
│   ├── utils/              # 工具函数
│   └── main.go            # 入口
├── frontend/               # React 前端
│   ├── src/
│   │   ├── assets/        # 静态资源
│   │   ├── components/    # UI 组件
│   │   ├── contexts/      # React Context
│   │   ├── i18n/          # 国际化
│   │   ├── lib/           # 工具库
│   │   ├── pages/         # 页面
│   │   └── utils/         # 工具函数
│   └── tests/             # 测试
├── Makefile               # 构建命令
├── AGENTS.md             # 开发规范
└── readme.md
```

## 快速开始

### 环境要求
- Go 1.25+
- Node.js 18+ / Bun
- PostgreSQL
- Redis

### 安装

```bash
# 后端
cd backend && go mod download

# 前端
cd frontend && npm install
# 或
cd frontend && bun install
```

### 开发

```bash
# 同时启动前后端
make dev

# 仅后端 (端口 8000)
make backend

# 仅前端 (端口 8001)
make frontend
```

### 构建

```bash
# 后端
make build-backend

# 前端
cd frontend && npm run build
```

## 测试

```bash
# 所有测试
make test

# 后端单元测试
make test-backend

# 前端单元测试
make test-frontend

# E2E 测试
make test-e2e
```

### 前端测试命令

```bash
cd frontend

# 单元测试
npm run test

# 测试 UI
npm run test:ui

# 测试覆盖率
npm run test:coverage

# 集成测试
npm run test:integration

# E2E 测试
npm run test:e2e
```

## 配置

### 后端
Viper 管理配置，支持 YAML/JSON/ENV。

### 前端
环境变量使用 `VITE_` 前缀：
- `VITE_API_BASE_URL` - API 基础 URL

## API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/register | 用户注册 |
| POST | /api/login | 用户登录 |

## Swagger

```bash
make swagger
```

访问：http://localhost:8000/swagger/index.html

## 功能

- 用户注册
- 用户登录
- JWT 认证
- CORS 支持
- 日志记录
- 错误处理
- 国际化

## 开发规范

详见 [AGENTS.md](./AGENTS.md)

## License

MIT
