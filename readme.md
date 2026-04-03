# Kaleidoscope

全栈应用，Go 后端 + React 前端，分层架构设计。

## 特性

### 后端
- **分层架构** - Controllers、Services、Models、Middleware 清晰分离
- **JWT 认证** - 完整的用户注册登录流程
- **Rate Limiting** - 基于 Redis 的请求限流
- **OpenTelemetry** - 分布式追踪和可观测性
- **结构化日志** - Zap 高性能日志
- **优雅关闭** - 信号处理和优雅关闭
- **Swagger 文档** - 自动生成 API 文档

### 前端
- **React 19** - 最新 React 特性
- **TypeScript** - strict 模式，完整类型安全
- **TailwindCSS 4** - 原子化 CSS
- **shadcn/ui** - 可访问性组件库
- **i18next** - 多语言支持
- **表单验证** - react-hook-form + zod
- **路由守卫** - React Router v7

### 测试
- **单元测试** - Vitest + React Testing Library
- **集成测试** - API 集成测试
- **E2E 测试** - Playwright
- **Mock** - MSW (Mock Service Worker)

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

## 文档

- [开发规范](./AGENTS.md) - 代码规范、架构设计、最佳实践
- [贡献指南](./CONTRIBUTING.md) - 开发流程、提交规范、PR 流程

## License

MIT
