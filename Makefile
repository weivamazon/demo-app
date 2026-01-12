.PHONY: build test run clean docker-build docker-run

VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags="-w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Build the application
build:
	@echo "Building demo-app..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o demo-app .

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run the application locally
run:
	@echo "Starting demo-app..."
	go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f demo-app coverage.out

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t demo-app:$(VERSION) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8000:8000 demo-app:$(VERSION)

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...

# All checks
check: fmt lint test
	@echo "All checks passed!"
