.PHONY: help test test-integration tidy

help: ## Show commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Download dependencies
	go mod tidy

test: ## Run unit tests only
	go test ./...

test-integration: ## Run integration tests with testcontainers (requires Docker)
	# -p 1: each OTP package boots its own Docker stack; parallel packages starve builds.
	go test -tags=integration -timeout=45m -p 1 -v ./tests/...

test-otp: ## Run OTP verify integration tests only
	go test -tags=integration -timeout=45m -p 1 -v ./tests/otp_verify/...

test-city-cache: ## Run city cache GET API integration tests
	go test -tags=integration -timeout=45m -p 1 -v ./tests/city_cache/...

test-category-cache: ## Run category cache GET API integration tests
	go test -tags=integration -timeout=45m -p 1 -v ./tests/category_cache/...
