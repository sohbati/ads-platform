.PHONY: help build run test clean swag docker-up docker-down

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building ads-platform..."
	@go build -o ads-platform main.go
	@echo "Build completed!"

run: ## Run the application
	@echo "Starting ads-platform..."
	@go run main.go

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f ads-platform
	@go clean

swag: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init
	@echo "Swagger docs generated!"

docker-up: ## Start MySQL with Docker Compose
	@echo "Starting MySQL..."
	@docker-compose up -d mysql
	@echo "MySQL started!"

docker-down: ## Stop all services
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped!"

dev: docker-up swag run ## Start development environment

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Dependencies installed!"

setup: install-deps docker-up swag ## Complete setup for development
	@echo "Setup completed! Run 'make run' to start the application."

