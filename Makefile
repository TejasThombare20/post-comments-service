# Post-Comments Service Makefile
# Development automation for Go application

# Variables
BINARY_NAME=post-comments-service
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/main.go
DOCKER_IMAGE=post-comments-service
DOCKER_TAG=latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOIMPORTS=goimports

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build clean test coverage deps tidy run dev fmt imports lint vet quality security docker-build docker-run migrate-up install-tools setup-hooks ci build-linux

# Default target
all: quality test build

# Help target
help: ## Show this help message
	@echo "$(BLUE)Post-Comments Service Development Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Building
build: ## Build the binary
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	@echo "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

build-linux: ## Build binary for Linux
	@echo "$(YELLOW)Building $(BINARY_UNIX) for Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)
	@echo "$(GREEN)Linux build completed: $(BINARY_UNIX)$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	@echo "$(GREEN)Clean completed$(NC)"

# Dependencies
deps: ## Install dependencies
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	$(GOGET) -d -v ./...
	@echo "$(GREEN)Dependencies installed$(NC)"

tidy: ## Tidy go modules
	@echo "$(YELLOW)Tidying go modules...$(NC)"
	$(GOMOD) tidy
	@echo "$(GREEN)Go modules tidied$(NC)"

# Running
run: ## Run the application
	@echo "$(YELLOW)Running $(BINARY_NAME)...$(NC)"
	$(GOCMD) run $(MAIN_PATH)

dev: ## Run with hot reload (requires air)
	@echo "$(YELLOW)Starting development server with hot reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)Air not found. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Falling back to regular run...$(NC)"; \
		$(MAKE) run; \
	fi

# Testing
test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)Tests completed$(NC)"

coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Code Quality
fmt: ## Format code with gofmt
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)Code formatted$(NC)"

imports: ## Organize imports with goimports
	@echo "$(YELLOW)Organizing imports...$(NC)"
	@if command -v $(GOIMPORTS) > /dev/null; then \
		$(GOIMPORTS) -w .; \
		echo "$(GREEN)Imports organized$(NC)"; \
	else \
		echo "$(RED)goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest$(NC)"; \
	fi

lint: ## Run golangci-lint
	@echo "$(YELLOW)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
		echo "$(GREEN)Linting completed$(NC)"; \
	else \
		echo "$(RED)golangci-lint not found. Install with: make install-tools$(NC)"; \
	fi

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GOCMD) vet ./...
	@echo "$(GREEN)Go vet completed$(NC)"

security: ## Run security scan with gosec
	@echo "$(YELLOW)Running security scan...$(NC)"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
		echo "$(GREEN)Security scan completed$(NC)"; \
	else \
		echo "$(RED)gosec not found. Install with: make install-tools$(NC)"; \
	fi

quality: fmt imports vet lint ## Run all quality checks
	@echo "$(GREEN)All quality checks completed$(NC)"

# Docker
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container
	@echo "$(YELLOW)Running Docker container...$(NC)"
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Database
migrate-up: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	@if [ -f "migrations/001_init.sql" ]; then \
		echo "$(BLUE)Migrations found, please run them manually or implement migration tool$(NC)"; \
	else \
		echo "$(RED)No migrations found$(NC)"; \
	fi

# Development Tools Setup
install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint > /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	@echo "Installing goimports..."
	@if ! command -v goimports > /dev/null; then \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "Installing gosec..."
	@if ! command -v gosec > /dev/null; then \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	@echo "Installing air (hot reload)..."
	@if ! command -v air > /dev/null; then \
		go install github.com/cosmtrek/air@latest; \
	fi
	@echo "Installing pre-commit..."
	@if ! command -v pre-commit > /dev/null; then \
		echo "$(YELLOW)Please install pre-commit manually:$(NC)"; \
		echo "  pip install pre-commit"; \
		echo "  or"; \
		echo "  brew install pre-commit"; \
	fi
	@echo "$(GREEN)Development tools installation completed$(NC)"

setup-hooks: ## Setup pre-commit hooks
	@echo "$(YELLOW)Setting up pre-commit hooks...$(NC)"
	@if command -v pre-commit > /dev/null; then \
		pre-commit install; \
		echo "$(GREEN)Pre-commit hooks installed$(NC)"; \
	else \
		echo "$(RED)pre-commit not found. Install with: pip install pre-commit$(NC)"; \
	fi

# CI/CD Simulation
ci: quality test build ## Simulate CI pipeline
	@echo "$(GREEN)CI pipeline simulation completed successfully$(NC)"

# Development workflow
dev-setup: install-tools setup-hooks tidy ## Complete development setup
	@echo "$(GREEN)Development environment setup completed$(NC)"

# Quick start
quick-start: dev-setup ## Quick start for new developers
	@echo "$(BLUE)Quick start completed! You can now:$(NC)"
	@echo "  $(GREEN)make run$(NC)     - Run the application"
	@echo "  $(GREEN)make dev$(NC)     - Run with hot reload"
	@echo "  $(GREEN)make test$(NC)    - Run tests"
	@echo "  $(GREEN)make quality$(NC) - Run quality checks" 