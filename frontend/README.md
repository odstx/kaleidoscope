# React + TypeScript + Vite

## Environment Configuration

- `.env` - Development environment
- `.env.test` - Test environment

### Port Configuration

| Environment | API Port | Frontend Port |
|-------------|----------|---------------|
| Development | 8000     | 8001          |
| Test        | 9000     | 9001          |

### Environment Variables

- `VITE_API_BASE_URL` - Backend API base URL
- `PLAYWRIGHT_TEST_BASE_URL` - Frontend URL for Playwright e2e tests (test env only)
- `CI` - CI environment flag (test env only)
