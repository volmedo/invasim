.PHONY: help lint tidy test check build clean

help: ## Show this help
	@echo "Available targets in this Makefile"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-15s\033[93m %s\n", $$1, $$2}'

lint: ## Lint code
	@bash scripts/lint.sh

tidy: ## Check go.mod and go.sum files are up to date
	@bash scripts/tidy.sh

test: ## Run unit tests
	@bash scripts/test.sh

check: lint tidy test ## Run all checks (lint, tidy, test)

build: ## Compile sources to get a binary
	@bash scripts/build.sh ./build/invasim

clean: ## Remove output directory
	@rm -r ./build
	