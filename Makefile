# Warp Launch Config Updater Makefile

BINARY_NAME=warp-config-updater
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: build clean install test

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go

# Install to /usr/local/bin
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Install to user's bin directory
install-user: build
	mkdir -p ~/bin
	mv $(BINARY_NAME) ~/bin/

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary for current platform"
	@echo "  build-all   - Build binaries for all platforms"
	@echo "  install     - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user- Install to ~/bin"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  help        - Show this help"
