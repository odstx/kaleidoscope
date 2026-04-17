.PHONY: dev backend frontend test test-backend test-frontend test-e2e swagger swag build-backend build build-frontend build-all _build-backend _build-frontend _copy-config image run macos deploy check check-frontend check-backend docker-up docker-down docker-build release install env env-new

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
else ifeq ($(ENV),uat)
	API_PORT := 8002
	FRONTEND_PORT := 8003
else ifeq ($(ENV),prod)
	API_PORT := 8000
	FRONTEND_PORT := 8001
else
	$(error Invalid ENV value. Use 'dev', 'test', 'uat' or 'prod')
endif

swagger:
	@echo "Generating Swagger documentation..."
	cd backend && swag init
	@echo "Swagger docs generated at backend/docs/"

swag: swagger

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
	@echo "Backend server will run on port $(API_PORT)"
	@echo "Backend worker will run for background tasks"
	@echo "Frontend will run on port $(FRONTEND_PORT)"
	@echo ""
	bash -c 'cd backend && go run . server --port $(API_PORT) & BACKEND_SERVER_PID=$$!; cd backend && go run . worker & BACKEND_WORKER_PID=$$!; cd $(PWD)/frontend && $(if $(filter test,$(ENV)),cp .env.test .env.local && bun run dev --port $(FRONTEND_PORT),bun run dev --port $(FRONTEND_PORT)) & FRONTEND_PID=$$!; echo "Backend Server PID: $$BACKEND_SERVER_PID"; echo "Backend Worker PID: $$BACKEND_WORKER_PID"; echo "Frontend PID: $$FRONTEND_PID"; trap "kill $$BACKEND_SERVER_PID $$BACKEND_WORKER_PID $$FRONTEND_PID 2>/dev/null; rm -f $(PWD)/frontend/.env.local" EXIT; wait $$BACKEND_SERVER_PID $$BACKEND_WORKER_PID $$FRONTEND_PID'

dev:
	@echo "Starting development environment..."
	@echo "Backend server will run on port 8000"
	@echo "Backend worker will run for background tasks"
	@echo "Frontend will run on port 8001"
	@echo ""
	bash -c 'cd backend && go run . server --port 8000 & BACKEND_SERVER_PID=$$!; cd backend && go run . worker & BACKEND_WORKER_PID=$$!; cd $(PWD)/frontend && bun run dev --port 8001 & FRONTEND_PID=$$!; echo "Backend Server PID: $$BACKEND_SERVER_PID"; echo "Backend Worker PID: $$BACKEND_WORKER_PID"; echo "Frontend PID: $$FRONTEND_PID"; trap "kill $$BACKEND_SERVER_PID $$BACKEND_WORKER_PID $$FRONTEND_PID 2>/dev/null" EXIT; wait $$BACKEND_SERVER_PID $$BACKEND_WORKER_PID $$FRONTEND_PID'

backend:
	cd backend && go run . server



image:
	@echo "Building Docker image for backend..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	cd backend && docker build -t kaleidoscope-backend:$(VERSION)-$(BUILD_ID) -f Dockerfile .
	@echo "Docker image built: kaleidoscope-backend:$(VERSION)-$(BUILD_ID)"

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
API_BASE_URL ?=
CONFIG_FILE ?=

build-all:
	@echo "Building frontend and backend to build directory..."
	@echo "Target: $(GOOS)/$(GOARCH)"
	@rm -rf build
	@mkdir -p build
	@$(MAKE) _build-backend
	@$(MAKE) _build-frontend
	@$(MAKE) _copy-config
	@echo "Build complete. Output in build/ directory"

 build: build-all

_build-backend:
	@echo "Building backend..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	cd backend && GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "\
		-X 'kaleidoscope/version.Version=$(VERSION)' \
		-X 'kaleidoscope/version.BuildID=$(BUILD_ID)' \
		-X 'kaleidoscope/version.BuildTime=$(BUILD_TIME)' \
		-X 'kaleidoscope/version.GitCommit=$(GIT_COMMIT)'" \
		-o ../build/kaleidoscope .
	@echo "Backend built successfully."

_build-frontend:
	@echo "Building frontend..."
	cd frontend && VITE_API_BASE_URL="$(API_BASE_URL)" bun run build
	@mv frontend/dist/* build/ 2>/dev/null || true
	@echo "Frontend built successfully."

 _copy-config:
	@echo "Copying config file..."
ifneq ($(CONFIG_FILE),)
	@cp $(CONFIG_FILE) build/config.yaml
else
	@cp backend/config/config.yaml build/config.yaml
endif
	@sed -i '' 's/environment: "development"/environment: "production"/' build/config.yaml

frontend:
	cd frontend && bun run dev

macos:
	@echo "Building macOS app from Swift package..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	cd swift && swift build -c release
	@mkdir -p swift/Kaleidoscope.app/Contents/MacOS
	@mkdir -p swift/Kaleidoscope.app/Contents/Resources
	@cp swift/.build/release/Kaleidoscope swift/Kaleidoscope.app/Contents/MacOS/
	@printf '<?xml version="1.0" encoding="UTF-8"?>\n<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">\n<plist version="1.0">\n<dict>\n<key>CFBundleDevelopmentRegion</key>\n<string>en</string>\n<key>CFBundleExecutable</key>\n<string>Kaleidoscope</string>\n<key>CFBundleIdentifier</key>\n<string>com.app.kaleidoscope</string>\n<key>CFBundleInfoDictionaryVersion</key>\n<string>6.0</string>\n<key>CFBundleName</key>\n<string>Kaleidoscope</string>\n<key>CFBundlePackageType</key>\n<string>APPL</string>\n<key>CFBundleShortVersionString</key>\n<string>$(VERSION)</string>\n<key>CFBundleVersion</key>\n<string>$(BUILD_ID)</string>\n<key>LSMinimumSystemVersion</key>\n<string>14.0</string>\n<key>NSHighResolutionCapable</key>\n<true/>\n</dict>\n</plist>' > swift/Kaleidoscope.app/Contents/Info.plist
	@echo "macOS app created at swift/Kaleidoscope.app"
	@echo "Killing existing Kaleidoscope..."
	@pkill -f Kaleidoscope 2>/dev/null || true
	@sleep 1
	@echo "Launching Kaleidoscope..."
	@open swift/Kaleidoscope.app

check:
	@echo "Running backend configuration check command..."
	cd backend && go run . check

deploy:
	@echo "Deploying to remote server (ENV=$(ENV))..."
	@bash deploy/deploy.sh $(ENV)
	@echo "Deployment to $(ENV) completed"

docker-up:
	@echo "Starting services with docker-compose..."
	cd backend && docker-compose up -d
	@echo "Services started! Backend available at http://localhost:8000"

docker-down:
	@echo "Stopping services with docker-compose..."
	cd backend && docker-compose down
	@echo "Services stopped."

docker-build:
	@echo "Building services with docker-compose..."
	cd backend && docker-compose build
	@echo "Services built."

release:
	@echo "Running goreleaser to create releases..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	cd backend && GITHUB_OWNER=odstx GITHUB_REPO=kaleidoscope goreleaser release --clean --timeout 10m

install:
	@echo "Installing dependencies for source deployment..."
	@echo "Installing Go module dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && bun install
	@echo "Installing Swift dependencies..."
	cd swift && swift package resolve
	@echo "Checking required tools..."
	@which psql >/dev/null 2>&1 || echo "Warning: psql not found. Install postgresql-client for database migrations."
	@which node >/dev/null 2>&1 || echo "Warning: node not found. Required for frontend build."
	@which bun >/dev/null 2>&1 || echo "Warning: bun not found. Required for frontend development."
	@echo "Installation complete!"

env-new:
ifndef ENV_NAME
	$(error ENV_NAME is required. Usage: make env-new ENV_NAME=prod)
endif
	@echo "Creating new deployment environment: $(ENV_NAME)"
	@mkdir -p deploy/$(ENV_NAME)
	@cp deploy/.env.example deploy/$(ENV_NAME)/.env
	@cp deploy/dev/config.yaml deploy/$(ENV_NAME)/config.yaml 2>/dev/null || true
	@echo "Created deploy/$(ENV_NAME)/.env and config.yaml"
	@echo "Please edit the files to configure the new environment"

env:
	@echo "Available deployment environments:"
	@ls -1 deploy/ | grep -E '^(dev|test|prod)$$' || echo "  (none found)"
	@echo ""
	@echo "To deploy: make deploy ENV=dev"
	@echo "To create new: make env-new ENV_NAME=prod"