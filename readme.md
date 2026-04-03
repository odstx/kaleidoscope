# Kaleidoscope

[![CI](https://github.com/odstx/kaleidoscope/workflows/CI/badge.svg)](https://github.com/odstx/kaleidoscope/actions/workflows/ci.yml)
[![Test Coverage](https://img.shields.io/badge/coverage-auto-brightgreen)](https://github.com/odstx/kaleidoscope/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

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

### 前端 (React)
- **React 19** - 最新 React 特性
- **TypeScript** - strict 模式，完整类型安全
- **TailwindCSS 4** - 原子化 CSS
- **shadcn/ui** - 可访问性组件库
- **i18next** - 多语言支持
- **表单验证** - react-hook-form + zod
- **路由守卫** - React Router v7

### 前端 (Swift)
- **SwiftUI** - 声明式 UI 框架
- **跨平台** - 支持 macOS 和 iOS
- **i18n** - 多语言支持 (英文/中文)
- **MVVM** - 清晰的架构分层
- **API Service** - 与后端完整集成
- **测试** - 单元测试和集成测试

### 测试
- **单元测试** - Vitest + React Testing Library
- **集成测试** - API 集成测试
- **E2E 测试** - Playwright
- **Mock** - MSW (Mock Service Worker)

## 文档

- [协作文档](./COLLABORATION.md) - 技术栈、快速开始、常用命令
- [开发规范](./AGENTS.md) - 代码规范、架构设计、最佳实践
- [贡献指南](./CONTRIBUTING.md) - 开发流程、提交规范、PR 流程

## License

MIT
