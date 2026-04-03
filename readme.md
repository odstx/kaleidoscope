# Kaleidoscope

全栈应用，Go 后端 + React 前端，分层架构设计。

## 特性

- **分层架构** - 对客层、接口层、应用层、数据层清晰分离
- **类型安全** - Go 强类型 + TypeScript strict 模式
- **现代化前端** - React 19 + Vite 8 + TailwindCSS 4
- **完整认证** - JWT 认证 + 用户注册登录
- **国际化** - i18next 多语言支持
- **API 文档** - Swagger 自动生成
- **测试完备** - 单元测试 + 集成测试 + E2E 测试

## 技术栈

**后端**: Go 1.25 + Gin + GORM + PostgreSQL + Redis

**前端**: React 19 + Vite 8 + TypeScript + TailwindCSS 4 + shadcn/ui

## 快速开始

```bash
# 安装依赖
cd backend && go mod download
cd frontend && bun install

# 启动开发服务器
make dev
```

后端运行在 http://localhost:8000，前端运行在 http://localhost:8001

## 文档

- [开发规范](./AGENTS.md) - 代码规范、架构设计、最佳实践
- [贡献指南](./CONTRIBUTING.md) - 开发流程、提交规范、PR 流程

## License

MIT
