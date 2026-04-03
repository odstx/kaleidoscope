.PHONY: dev backend frontend test test-backend test-frontend test-e2e swagger build-backend run

VERSION := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "dev")
BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

ENV ?= dev

ifeq ($(ENV),dev)
	API_PORT := 8000
	FRONTEND_PORT := 8001
else ifeq ($(ENV),test)
	API_PORT := 9000
	FRONTEND_PORT := 9001
else
	$(error Invalid ENV value. Use 'dev' or 'test')
endif

swagger:
	@echo "Generating Swagger documentation..."
	cd backend && swag init
	@echo "Swagger docs generated at backend/docs/"

test:
	@echo "Running all tests..."
	@echo ""
	@echo "=== Backend Unit Tests ==="
	$(MAKE) test-backend
	@echo ""
	@echo "=== Frontend Unit Tests ==="
	$(MAKE) test-frontend
	@echo ""
	@echo "=== E2E Tests ==="
	$(MAKE) test-e2e
	@echo ""
	@echo "All tests completed!"

test-backend:
	cd backend && go test ./... -v

test-frontend:
	cd frontend && bun run test --run

test-e2e:
	cd frontend && bun run test:e2e

run:
	@echo "Starting $(ENV) environment..."
	@echo "Backend will run on port $(API_PORT)"
	@echo "Frontend will run on port $(FRONTEND_PORT)"
	@echo ""
	bash -c 'cd backend && go run . server --port $(API_PORT) & BACKEND_PID=$$!; cd $(PWD)/frontend && $(if $(filter test,$(ENV)),cp .env.test .env.local && bun run dev --port $(FRONTEND_PORT),bun run dev --port $(FRONTEND_PORT)) & FRONTEND_PID=$$!; echo "Backend PID: $$BACKEND_PID"; echo "Frontend PID: $$FRONTEND_PID"; trap "kill $$BACKEND_PID $$FRONTEND_PID 2>/dev/null; rm -f $(PWD)/frontend/.env.local" EXIT; wait $$BACKEND_PID $$FRONTEND_PID'

dev:
	@echo "Starting development environment..."
	@echo "Backend will run on port 8000"
	@echo "Frontend will run on port 8001"
	@echo ""
	bash -c 'cd backend && go run . server --port 8000 & BACKEND_PID=$$!; cd $(PWD)/frontend && bun run dev --port 8001 & FRONTEND_PID=$$!; echo "Backend PID: $$BACKEND_PID"; echo "Frontend PID: $$FRONTEND_PID"; trap "kill $$BACKEND_PID $$FRONTEND_PID 2>/dev/null" EXIT; wait $$BACKEND_PID $$FRONTEND_PID'

backend:
	cd backend && go run . server

build-backend:
	@echo "Building backend with version info..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	cd backend && go build -ldflags "\
		-X 'kaleidoscope/version.Version=$(VERSION)' \
		-X 'kaleidoscope/version.BuildID=$(BUILD_ID)' \
		-X 'kaleidoscope/version.BuildTime=$(BUILD_TIME)' \
		-X 'kaleidoscope/version.GitCommit=$(GIT_COMMIT)'" \
		-o kaleidoscope .

frontend:
	cd frontend && bun run dev