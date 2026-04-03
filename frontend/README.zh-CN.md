# React + TypeScript + Vite

## 环境配置

- `.env` - 开发环境
- `.env.test` - 测试环境

### 端口配置

| 环境 | API端口 | 前端端口 |
|------|---------|----------|
| 开发 | 8000    | 8001     |
| 测试 | 9000    | 9001     |

### 环境变量

- `VITE_API_BASE_URL` - 后端 API 基础 URL
- `PLAYWRIGHT_TEST_BASE_URL` - Playwright e2e 测试的前端 URL（仅测试环境）
- `CI` - CI 环境标志（仅测试环境）
