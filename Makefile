NAME ?= gogobosh
GO_SOURCES = $(shell find . -type f -name '*.go')
GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT_VERSION := $(shell golangci-lint --version 2>/dev/null)

.PHONY: all
all: build lint test ## Runs build, lint and test

.PHONY: clean
clean: ## Clean testcache and delete build output
	go clean -testcache
	@rm -f coverage.html

.PHONY: build
build:
	go build ./...

.PHONY: test
test: ## Run the unit tests
	go test -short ./...

.PHONY: test-integration
test-integration: ## Run integration tests only
	go test -timeout 15m -v -tags integration ./integration_test.go

.PHONY: test-all
test-all: test test-integration ## Run integration & unit tests

.PHONY: coverage
coverage: ## Run the tests with coverage and race detection
	go test -v --race -coverprofile=c.out -covermode=atomic ./...

.PHONY: report
report: ## Show coverage in an html report
	go tool cover -html=c.out -o coverage.html

.PHONY: lint
lint: ## Validate style and syntax
ifdef GOLANGCI_LINT_VERSION
	golangci-lint run
else
	@echo "Installing latest golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	@echo "[OK] golangci-lint installed"
	./bin/golangci-lint run
endif

.PHONY: tidy
tidy: ## Remove unused dependencies
	go mod tidy

.PHONY: list
list: ## Print the current module's dependencies.
	go list -m all

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
