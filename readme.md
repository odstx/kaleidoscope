# Kaleidoscope

[![CI](https://github.com/odstx/kaleidoscope/workflows/CI/badge.svg)](https://github.com/odstx/kaleidoscope/actions/workflows/ci.yml)
[![Test Coverage](https://img.shields.io/badge/coverage-auto-brightgreen)](https://github.com/odstx/kaleidoscope/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Harness 驱动的全栈应用。

## 技术栈

### 后端 (Go 1.25)
Gin + GORM + PostgreSQL + Redis，分层架构（Controllers、Services、Models、Middleware），JWT 认证，Rate Limiting，OpenTelemetry，Zap 日志，Swagger 文档

### 前端 (React 19 + Vite 8)
TypeScript strict 模式，TailwindCSS 4，shadcn/ui，i18next 多语言，react-hook-form + zod 表单验证，React Router v7

### 前端 (SwiftUI)
跨平台 macOS/iOS，MVVM 架构，多语言支持，完整 API 集成

### 前端 (Kotlin + Compose)
Android 应用，Jetpack Compose UI，Hilt 依赖注入，Retrofit + OkHttp 网络层，DataStore 本地存储，MVVM 架构

### 测试
Vitest + React Testing Library 单元测试，Playwright E2E 测试，MSW Mock

## 文档

- [协作文档](./COLLABORATION.md) - 快速开始、常用命令、部署
- [开发规范](./AGENTS.md) - 代码规范、架构设计
- [贡献指南](./CONTRIBUTING.md) - 开发流程、PR 规范
- [安装部署](./install.md) - 环境配置、部署流程

## License

MIT
