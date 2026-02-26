.PHONY: all build dev clean install-deps

# Default target
all: build

# Install dependencies
install-deps:
	cd frontend && npm install
	go mod tidy

# Development mode with hot reload
dev:
	wails dev

# Build for production
build:
	wails build

# Build for Windows
build-windows:
	wails build -platform windows/amd64

# Build for Linux
build-linux:
	wails build -platform linux/amd64

# Build for macOS
build-darwin:
	wails build -platform darwin/amd64
	wails build -platform darwin/arm64

# Build all platforms
build-all:
	wails build -platform windows/amd64,linux/amd64,darwin/amd64,darwin/arm64

# Clean build artifacts
clean:
	rm -rf build/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...
	cd frontend && npm run lint

# Generate Wails bindings
generate:
	wails generate module
