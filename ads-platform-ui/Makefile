.PHONY: help build run test clean tidy setup

APP_NAME=ads-platform-ui

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) ./cmd/server
	@echo "Build completed!"

run: ## Run the application
	@echo "Starting $(APP_NAME)..."
	@go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(APP_NAME)
	@go clean

tidy: ## Tidy Go modules
	@go mod tidy

setup: tidy ## Complete setup
	@echo "Setup completed! Run 'make run' to start."
