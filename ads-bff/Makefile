.PHONY: help run tidy test

help: ## Show commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Download dependencies
	go mod tidy

run: ## Run locally
	go run ./cmd/server

test: ## Run unit tests
	go test ./...
