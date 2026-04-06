.PHONY: dev backend frontend test test-backend test-frontend test-e2e swagger build-backend build run macos deploy

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

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
API_BASE_URL ?= 

build:
	@echo "Building frontend and backend to build directory..."
	@echo "Target: $(GOOS)/$(GOARCH)"
	@rm -rf build
	@mkdir -p build
	@echo "Building backend..."
	cd backend && GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "\
		-X 'kaleidoscope/version.Version=$(VERSION)' \
		-X 'kaleidoscope/version.BuildID=$(BUILD_ID)' \
		-X 'kaleidoscope/version.BuildTime=$(BUILD_TIME)' \
		-X 'kaleidoscope/version.GitCommit=$(GIT_COMMIT)'" \
		-o ../build/kaleidoscope .
	@echo "Building frontend..."
	cd frontend && VITE_API_BASE_URL="$(API_BASE_URL)" bun run build
	@echo "Copying config file..."
	@cp backend/config/config.yaml build/config.yaml
	@sed -i '' 's/environment: "development"/environment: "production"/' build/config.yaml
	@echo "Build complete. Output in build/ directory"

frontend:
	cd frontend && bun run dev

macos:
	@echo "Building macOS app from Swift package..."
	@echo "Version: $(VERSION)"
	@echo "Build ID: $(BUILD_ID)"
	cd swift && swift build -c release
	@mkdir -p swift/FrontendApp.app/Contents/MacOS
	@mkdir -p swift/FrontendApp.app/Contents/Resources
	@cp swift/.build/release/FrontendApp swift/FrontendApp.app/Contents/MacOS/
	@printf '<?xml version="1.0" encoding="UTF-8"?>\n<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">\n<plist version="1.0">\n<dict>\n<key>CFBundleDevelopmentRegion</key>\n<string>en</string>\n<key>CFBundleExecutable</key>\n<string>FrontendApp</string>\n<key>CFBundleIdentifier</key>\n<string>com.app.frontend</string>\n<key>CFBundleInfoDictionaryVersion</key>\n<string>6.0</string>\n<key>CFBundleName</key>\n<string>FrontendApp</string>\n<key>CFBundlePackageType</key>\n<string>APPL</string>\n<key>CFBundleShortVersionString</key>\n<string>$(VERSION)</string>\n<key>CFBundleVersion</key>\n<string>$(BUILD_ID)</string>\n<key>LSMinimumSystemVersion</key>\n<string>14.0</string>\n<key>NSHighResolutionCapable</key>\n<true/>\n</dict>\n</plist>' > swift/FrontendApp.app/Contents/Info.plist
	@echo "macOS app created at swift/FrontendApp.app"
	@echo "Killing existing FrontendApp..."
	@pkill -f FrontendApp 2>/dev/null || true
	@sleep 1
	@echo "Launching FrontendApp..."
	@open swift/FrontendApp.app

deploy:
	@echo "Deploying to remote server..."
	@if [ ! -f deploy/.env ]; then \
		echo "Error: deploy/.env not found"; \
		echo "Please configure deployment:"; \
		echo "  1. cp deploy/.env.example deploy/.env"; \
		echo "  2. Edit deploy/.env with your server details"; \
		exit 1; \
	fi
	@bash deploy/deploy.sh