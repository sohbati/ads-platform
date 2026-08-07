.PHONY: help test test-integration tidy

help: ## Show commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Download dependencies
	go mod tidy

test: ## Run unit tests only
	go test ./...

test-integration: ## Run integration tests with testcontainers (requires Docker)
	go test -tags=integration -timeout=20m -v ./tests/...
