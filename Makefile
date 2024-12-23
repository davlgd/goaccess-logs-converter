.PHONY: build clean install lint test build-all help

BINARY_NAME=goaccess-logs-converter
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.name=${BINARY_NAME} -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
GOFLAGS=-trimpath

default: build

build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${GOFLAGS} ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/${BINARY_NAME}

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@go clean

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

build-all: clean
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build ${GOFLAGS} ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/${BINARY_NAME}
	@GOOS=darwin GOARCH=amd64 go build ${GOFLAGS} ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd/${BINARY_NAME}
	@GOOS=windows GOARCH=amd64 go build ${GOFLAGS} ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd/${BINARY_NAME}

install:
	@echo "Installing ${BINARY_NAME}..."
	@go install ${GOFLAGS} ${LDFLAGS} ./cmd/${BINARY_NAME}

help:
	@echo "Available commands:"
	@echo "  build       - Build optimized binary"
	@echo "  build-all   - Build optimized binaries for multiple platforms"
	@echo "  clean       - Clean up build artifacts"
	@echo "  install     - Install binary in go bin"
	@echo "  lint        - Run the linter"
	@echo "  test        - Run tests"
