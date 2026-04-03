# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-04-03

### Added
- Initial project structure with Go backend and React frontend
- Navbar component showing APP_NAME from environment variables
- JWT authentication system
- i18next internationalization support
- Swagger API documentation
- Comprehensive test setup (unit, integration, E2E)

### Changed
- Migrated from npm to Bun runtime for frontend
- Simplified README to focus on project features
- Updated AGENTS.md to use Bun commands

### Technical Details
- Frontend: React 19 + Vite 8 + TypeScript + TailwindCSS 4 + shadcn/ui
- Backend: Go 1.25 + Gin + GORM + PostgreSQL + Redis
- Testing: Vitest + React Testing Library + Playwright + MSW

---

## Version History

| Date | Version | Description |
|------|---------|-------------|
| 2026-04-03 | 0.1.0 | Initial release with Bun migration |

---

For detailed daily changes, see individual daily logs in `/changelog/daily/` directory.
