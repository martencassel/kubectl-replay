.PHONY: help
help: ## Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: build
build: ## Build the kubectl-replay binary
	go build -o bin/kubectl-replay main.go

.PHONY: release
release: ## Create release tarball with version from plugin.yaml
	@./hack/release.sh

.DEFAULT_GOAL := help
