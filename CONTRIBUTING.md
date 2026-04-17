# Contributing Guide

Thank you for your interest in contributing! This guide will help you understand how to participate in the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Environment](#development-environment)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Commit Standards](#commit-standards)
- [Branch Strategy](#branch-strategy)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

- Respect all contributors
- Maintain professional and constructive discussions
- Accept constructive criticism
- Focus on what benefits the community

## Development Environment

### System Requirements

- Go 1.25+
- Bun 1.3+
- PostgreSQL 14+
- Redis 7+
- Git

### Installation Steps

1. Fork and clone the repository

```bash
git clone https://github.com/YOUR_USERNAME/kaleidoscope.git
cd kaleidoscope
```

2. Install backend dependencies

```bash
cd backend
go mod download
```

3. Install frontend dependencies

```bash
cd frontend
bun install
```

4. Configure environment variables

```bash
# Backend
cp backend/.env.example backend/.env

# Frontend
cp frontend/.env.example frontend/.env.local
```

5. Start development server

```bash
make dev
```

## Development Workflow

### 1. Create a feature branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make your changes

- Follow code standards in [AGENTS.md](./AGENTS.md)
- Write clear code and comments
- Add necessary tests

### 3. Test locally

```bash
# Run all tests
make test

# Or run separately
make test-backend  # Backend tests
make test-frontend # Frontend tests
make test-e2e      # E2E tests
```

### 4. Run code checks

```bash
# Backend
cd backend && go fmt ./... && go vet ./...

# Frontend
cd frontend && bun run lint
```

### 5. Commit changes

```bash
git add .
git commit -m "feat: add your feature description"
```

### 6. Push and create PR

```bash
git push origin feature/your-feature-name
```

## Code Standards

### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Format code with `gofmt`
- Check code with `go vet`
- Add appropriate comments and docstrings
- Handle errors explicitly, don't ignore them

### Frontend (TypeScript/React)

See [AGENTS.md](./AGENTS.md) for details.

Key points:
- Use TypeScript strict mode
- Use functional components with hooks
- Use `@/` alias for imports
- Follow shadcn/ui component standards
- Use react-hook-form + zod for all forms

## Commit Standards

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation update
- `style`: Code formatting (no functional change)
- `refactor`: Refactoring
- `perf`: Performance optimization
- `test`: Test related
- `chore`: Build/tool related
- `ci`: CI/CD related

### Scope

- `backend`: Backend related
- `frontend`: Frontend related
- `api`: API related
- `ui`: UI component related
- `auth`: Authentication related
- `i18n`: i18n related

### Examples

```bash
feat(frontend): add user profile page
fix(backend): resolve database connection timeout issue
docs: update installation instructions
refactor(api): simplify authentication middleware
test(frontend): add unit tests for LoginForm component
```

## Branch Strategy

### Branch Naming

- `main`: Main branch, production code
- `develop`: Development branch (if applicable)
- `feature/description`: New feature branch
- `fix/description`: Bug fix branch
- `refactor/description`: Refactoring branch
- `docs/description`: Documentation update branch
- `test/description`: Test related branch

### Branch Flow

```
main (production)
  ↑
  └── feature/xxx → PR → merge
  └── fix/xxx → PR → merge
```

1. Create feature branch from `main`
2. Complete development and testing
3. Create Pull Request
4. Pass code review and CI checks
5. Merge to `main`

## Pull Request Process

### Pre-PR Checklist

- [ ] Code follows project standards
- [ ] All tests pass
- [ ] New features have corresponding tests
- [ ] Documentation updated (if needed)
- [ ] Commit message follows standards
- [ ] Branch created from latest `main`

### PR Title Format

```
<type>(<scope>): <description>
```

Example: `feat(frontend): add dark mode support`

### PR Description Template

```markdown
## Change Type
- [ ] New feature (feat)
- [ ] Bug fix (fix)
- [ ] Refactor (refactor)
- [ ] Documentation update (docs)
- [ ] Other

## Description
<!-- Describe what and why -->

## Related Issue
<!-- Link related issue, e.g: Closes #123 -->

## Testing
<!-- How to test this change -->

## Screenshots
<!-- If UI changes, provide screenshots -->

## Checklist
- [ ] Code follows standards
- [ ] Tests pass
- [ ] Documentation updated
```

### Code Review

- Each PR requires at least 1 approval
- Reviewers should check:
  - Code quality and standards
  - Test coverage
  - Potential issues and improvements
  - Documentation completeness

## Testing Requirements

### Backend Tests

```bash
# Unit tests
cd backend && go test ./... -v

# Test coverage
cd backend && go test ./... -cover

# Specific package test
cd backend && go test ./services -v
```

### Frontend Tests

```bash
cd frontend

# Unit tests
bun run test

# UI tests
bun run test:ui

# Test coverage
bun run test:coverage

# Integration tests
bun run test:integration

# E2E tests
bun run test:e2e
```

### Testing Standards

- All new features must have corresponding tests
- Bug fixes should include regression tests
- Unit test coverage > 80%
- Key paths need E2E tests
- Test code should also follow code standards

## Issue Reporting

### Bug Reports

Include when creating an issue:

1. **Description**: Clear description of the problem
2. **Steps to Reproduce**: How to reproduce
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happened
5. **Environment Info**:
   - OS: [e.g. macOS 14]
   - Go version: [e.g. 1.25]
   - Bun version: [e.g. 1.3.11]
6. **Screenshots**: If applicable
7. **Logs**: Relevant error logs

### Feature Requests

Include:

1. **Feature Description**: What you want
2. **Use Case**: Why you need this
3. **Implementation Suggestions**: Optional implementation ideas
4. **Alternatives**: Other solutions considered

## Documentation

### Updating Documentation

- API changes need Swagger doc updates
- New features need README.md updates
- Config changes need environment variable documentation
- Architecture changes need architecture diagram updates

### Generate Swagger Docs

```bash
make swagger
```

Visit: http://localhost:8000/swagger/index.html

## Getting Help

- See [README.md](./README.md) for project overview
- See [AGENTS.md](./AGENTS.md) for development standards
- Ask questions in Issues
- Check existing Issues and PRs

## License

This project uses MIT license. Contributed code will use the same license.