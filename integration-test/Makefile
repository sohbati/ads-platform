.PHONY: help test test-integration tidy

help: ## Show commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Download dependencies
	go mod tidy

test: ## Run unit tests only
	go test ./...

test-integration: ## Run integration tests with testcontainers (requires Docker)
	# -p 1: suites share stacks internally; avoid parallel packages competing for Docker.
	go test -tags=integration -timeout=45m -p 1 -v ./tests/...

test-otp: ## Run OTP verify suite (one shared stack)
	go test -tags=integration -timeout=45m -v ./tests/otp_verify/

test-cache-catalog: ## Run city+category+attr-schema cache suite (one shared CDN+cache stack)
	go test -tags=integration -timeout=45m -v ./tests/cache_catalog/

test-ads-api: ## Run POST /ads suite (postgres+cdn+cache+minio+back)
	go test -tags=integration -timeout=45m -v ./tests/ads_api/
