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
```
