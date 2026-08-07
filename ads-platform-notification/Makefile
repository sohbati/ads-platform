.PHONY: help build run test clean install-deps setup

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building ads-platform-notification..."
	@go build -o ads-platform-notification ./cmd/server
	@echo "Build completed!"

run: ## Run the application
	@echo "Starting ads-platform-notification..."
	@go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f ads-platform-notification
	@go clean

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Dependencies installed!"

setup: install-deps ## Complete setup for development
	@echo "Setup completed! Run 'make run' to start the application."
