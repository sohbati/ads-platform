.PHONY: help build run test clean swag install-deps setup

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building message-broker..."
	@go build -o message-broker ./cmd/server
	@echo "Build completed!"

run: ## Run the application
	@echo "Starting message-broker..."
	@go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f message-broker
	@go clean

swag: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go
	@echo "Swagger docs generated!"

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Dependencies installed!"

setup: install-deps ## Complete setup for development
	@echo "Setup completed! Run 'make run' to start the application."
