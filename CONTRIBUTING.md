# 贡献指南

感谢您对本项目的关注！本文档将帮助您了解如何参与项目开发。

## 目录

- [行为准则](#行为准则)
- [开发环境设置](#开发环境设置)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [分支策略](#分支策略)
- [Pull Request 流程](#pull-request-流程)
- [测试要求](#测试要求)
- [问题反馈](#问题反馈)

## 行为准则

- 尊重所有贡献者
- 保持专业和建设性的讨论
- 接受建设性批评
- 关注对社区最有利的事情

## 开发环境设置

### 系统要求

- Go 1.25+
- Bun 1.3+
- PostgreSQL 14+
- Redis 7+
- Git

### 安装步骤

1. Fork 并克隆仓库

```bash
git clone https://github.com/YOUR_USERNAME/kaleidoscope.git
cd kaleidoscope
```

2. 安装后端依赖

```bash
cd backend
go mod download
```

3. 安装前端依赖

```bash
cd frontend
bun install
```

4. 配置环境变量

```bash
# 后端
cp backend/.env.example backend/.env

# 前端
cp frontend/.env.example frontend/.env.local
```

5. 启动开发服务器

```bash
make dev
```

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 进行开发

- 遵循 [AGENTS.md](./AGENTS.md) 中的代码规范
- 编写清晰的代码和注释
- 添加必要的测试

### 3. 本地测试

```bash
# 运行所有测试
make test

# 或分别运行
make test-backend  # 后端测试
make test-frontend # 前端测试
make test-e2e      # E2E 测试
```

### 4. 代码检查

```bash
# 后端
cd backend && go fmt ./... && go vet ./...

# 前端
cd frontend && bun run lint
```

### 5. 提交更改

```bash
git add .
git commit -m "feat: add your feature description"
```

### 6. 推送并创建 PR

```bash
git push origin feature/your-feature-name
```

## 代码规范

### 后端 (Go)

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 规范
- 使用 `gofmt` 格式化代码
- 使用 `go vet` 检查代码
- 添加适当的注释和文档字符串
- 错误处理要明确，不要忽略错误

### 前端 (TypeScript/React)

详见 [AGENTS.md](./AGENTS.md)

关键点：
- 使用 TypeScript strict 模式
- 组件使用函数式组件 + hooks
- 使用 `@/` 别名进行导入
- 遵循 shadcn/ui 组件规范
- 所有表单使用 react-hook-form + zod

## 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具相关
- `ci`: CI/CD 相关

### Scope 范围

- `backend`: 后端相关
- `frontend`: 前端相关
- `api`: API 相关
- `ui`: UI 组件相关
- `auth`: 认证相关
- `i18n`: 国际化相关

### 示例

```bash
feat(frontend): add user profile page
fix(backend): resolve database connection timeout issue
docs: update installation instructions
refactor(api): simplify authentication middleware
test(frontend): add unit tests for LoginForm component
```

## 分支策略

### 分支命名

- `main`: 主分支，生产环境代码
- `develop`: 开发分支（如果有）
- `feature/description`: 新功能分支
- `fix/description`: Bug 修复分支
- `refactor/description`: 重构分支
- `docs/description`: 文档更新分支
- `test/description`: 测试相关分支

### 分支流程

```
main (生产)
  ↑
  └── feature/xxx → PR → merge
  └── fix/xxx → PR → merge
```

1. 从 `main` 创建功能分支
2. 完成开发和测试
3. 创建 Pull Request
4. 通过代码审查和 CI 检查
5. 合并到 `main`

## Pull Request 流程

### 创建 PR 前检查清单

- [ ] 代码遵循项目规范
- [ ] 所有测试通过
- [ ] 新功能有对应测试
- [ ] 文档已更新（如需要）
- [ ] 提交信息符合规范
- [ ] 分支从最新 `main` 创建

### PR 标题格式

```
<type>(<scope>): <description>
```

示例：`feat(frontend): add dark mode support`

### PR 描述模板

```markdown
## 变更类型
- [ ] 新功能 (feat)
- [ ] Bug 修复 (fix)
- [ ] 重构 (refactor)
- [ ] 文档更新 (docs)
- [ ] 其他

## 变更说明
<!-- 描述本次变更的内容和原因 -->

## 相关 Issue
<!-- 关联的 Issue 编号，如: Closes #123 -->

## 测试说明
<!-- 如何测试本次变更 -->

## 截图
<!-- 如有 UI 变更，提供截图 -->

## 检查清单
- [ ] 代码遵循规范
- [ ] 测试通过
- [ ] 文档已更新
```

### 代码审查

- 每个 PR 需要至少 1 个审查批准
- 审查者应检查：
  - 代码质量和规范
  - 测试覆盖率
  - 潜在问题和改进建议
  - 文档完整性

## 测试要求

### 后端测试

```bash
# 单元测试
cd backend && go test ./... -v

# 测试覆盖率
cd backend && go test ./... -cover

# 特定包测试
cd backend && go test ./services -v
```

### 前端测试

```bash
cd frontend

# 单元测试
bun run test

# 测试 UI
bun run test:ui

# 测试覆盖率
bun run test:coverage

# 集成测试
bun run test:integration

# E2E 测试
bun run test:e2e
```

### 测试规范

- 所有新功能必须有对应测试
- Bug 修复应包含回归测试
- 单元测试覆盖率 > 80%
- 关键路径需要 E2E 测试
- 测试代码也要遵循代码规范

## 问题反馈

### Bug 报告

创建 Issue 时请包含：

1. **问题描述**: 清晰描述问题
2. **复现步骤**: 如何复现问题
3. **期望行为**: 期望发生什么
4. **实际行为**: 实际发生了什么
5. **环境信息**: 
   - OS: [e.g. macOS 14]
   - Go 版本: [e.g. 1.25]
   - Bun 版本: [e.g. 1.3.11]
6. **截图**: 如适用
7. **日志**: 相关错误日志

### 功能请求

请包含：

1. **功能描述**: 想要什么功能
2. **使用场景**: 为什么需要这个功能
3. **实现建议**: 可选的实现思路
4. **替代方案**: 考虑过的其他方案

## 文档

### 更新文档

- API 变更需更新 Swagger 文档
- 新功能需更新 README.md
- 配置变更需更新环境变量说明
- 架构变更需更新架构图

### 生成 Swagger 文档

```bash
make swagger
```

访问：http://localhost:8000/swagger/index.html

## 获取帮助

- 查看 [README.md](./README.md) 了解项目概况
- 查看 [AGENTS.md](./AGENTS.md) 了解开发规范
- 在 Issue 中提问
- 查看 existing Issues 和 PRs

## 许可证

本项目采用 MIT 许可证。贡献的代码将采用相同许可证。
