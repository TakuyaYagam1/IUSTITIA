.DEFAULT_GOAL := help
SHELL         := bash
.SHELLFLAGS   := -c

# ------------------------------------------------------------------------------
# Pinned tool versions
# ------------------------------------------------------------------------------
OAPI_CODEGEN_VERSION  ?= v2.5.1
SQLC_VERSION          ?= v1.30.0
GOOSE_VERSION         ?= v3.27.0
WIRE_VERSION          ?= v0.7.0
GOLANGCI_LINT_VERSION ?= v2.11.4
GOVULNCHECK_VERSION   ?= v1.2.0
REDOCLY_VERSION       ?= 2.29.0

# ------------------------------------------------------------------------------
# Paths
# ------------------------------------------------------------------------------
OPENAPI_SRC    := internal/openapi/openapi.yml
OPENAPI_BUNDLE := internal/openapi/openapi.yaml
DB_PATH        ?= ./iustitia.db
MIGRATIONS_DIR := migrations
LINT_TIMEOUT   ?= 5m
TEST_TIMEOUT   ?= 5m

# ------------------------------------------------------------------------------
# HELP
# ------------------------------------------------------------------------------
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup

.PHONY: install-tools

install-tools: ## Install pinned dev tools (oapi-codegen, sqlc, goose, wire, golangci-lint, govulncheck)
	@echo "Installing pinned tools..."
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	@go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@go install github.com/google/wire/cmd/wire@$(WIRE_VERSION)
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	@echo "Tools installed. Note: redocly-cli вызывается через 'npx --yes @redocly/cli@$(REDOCLY_VERSION)' - ставить отдельно не нужно."

##@ Code generation

.PHONY: openapi-bundle openapi oapi-codegen sqlc sqlc-verify validate-openapi wire generate

openapi-bundle: ## Bundle split OpenAPI spec (openapi.yml + routes/*) into single openapi.yaml
	@echo "Bundling OpenAPI spec..."
	@bash -c 'npx --yes @redocly/cli@$(REDOCLY_VERSION) bundle $(OPENAPI_SRC) -o $(OPENAPI_BUNDLE) --ext yaml'
	@echo "Bundled to $(OPENAPI_BUNDLE)"

oapi-codegen: openapi-bundle ## Generate Go code from bundled OpenAPI (types + server + spec)
	@echo "Generating OpenAPI Go code..."
	@oapi-codegen -config codegen/oapi-codegen-types.yml  $(OPENAPI_BUNDLE)
	@oapi-codegen -config codegen/oapi-codegen-server.yml $(OPENAPI_BUNDLE)
	@oapi-codegen -config codegen/oapi-codegen-spec.yml   $(OPENAPI_BUNDLE)
	@rm -f $(OPENAPI_BUNDLE)
	@echo "OpenAPI code generated (types.gen.go, server.gen.go, spec.gen.go)"

openapi: oapi-codegen ## Alias for oapi-codegen (bundle + generate)

validate-openapi: openapi-bundle ## Lint bundled OpenAPI spec (redocly)
	@echo "Validating OpenAPI..."
	@bash -c 'npx --yes @redocly/cli@$(REDOCLY_VERSION) lint $(OPENAPI_BUNDLE)'
	@rm -f $(OPENAPI_BUNDLE)
	@echo "OpenAPI validation OK"

sqlc-verify: ## Verify sqlc schema and queries compile without generating code
	@echo "Verifying sqlc..."
	@sqlc -f codegen/sqlc.yaml compile
	@echo "sqlc verification OK"

sqlc: sqlc-verify ## Generate sqlc Go code from queries and schema
	@echo "Generating sqlc code..."
	@sqlc -f codegen/sqlc.yaml generate
	@echo "sqlc code generated"

wire: ## Generate Wire DI code
	@echo "Generating Wire code..."
	@cd internal/wire && go generate .
	@echo "Wire code generated"

generate: openapi sqlc wire ## Run full code generation pipeline (openapi + sqlc + wire)
	@echo "All code generation completed"

##@ Database

.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up: ## Apply all pending migrations (DB_PATH=$(DB_PATH))
	@echo "Applying migrations..."
	@goose -dir $(MIGRATIONS_DIR) sqlite "$(DB_PATH)" up

migrate-down: ## Revert last migration batch
	@echo "Reverting last migration..."
	@goose -dir $(MIGRATIONS_DIR) sqlite "$(DB_PATH)" down

migrate-status: ## Show migration status
	@goose -dir $(MIGRATIONS_DIR) sqlite "$(DB_PATH)" status

migrate-create: ## Create new migration (usage: make migrate-create name=add_something)
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

##@ Build & Run

.PHONY: build run clean

build: ## Build the backend binary (run 'make generate' first if spec/sql changed)
	@echo "Building binary..."
	@go build -o bin/iustitia ./cmd/app
	@echo "Built: bin/iustitia"

run: ## Run the application locally (uses DB_PATH=$(DB_PATH))
	@echo "Running iustitia backend..."
	@DB_PATH=$(DB_PATH) go run ./cmd/app

clean: ## Remove build artefacts and generated temp files
	@echo "Cleaning..."
	@rm -rf bin
	@rm -f coverage.out coverage.html
	@rm -f $(OPENAPI_BUNDLE)
	@echo "Clean completed"

##@ Tests

.PHONY: test test-race

test: ## Run all tests (no race detector)
	@echo "Running tests..."
	@go test -parallel 5 -timeout $(TEST_TIMEOUT) ./...

test-race: ## Run all tests with race detector (slower)
	@echo "Running tests (race)..."
	@go test -race -parallel 5 -timeout $(TEST_TIMEOUT) ./...

##@ Quality

.PHONY: fmt lint lint-fix mod-tidy deps-audit

fmt: ## Format code with gofmt
	@echo "Formatting..."
	@gofmt -w -s .
	@echo "Done"

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run -v --timeout $(LINT_TIMEOUT) ./...

lint-fix: ## Run golangci-lint with --fix
	@echo "Running linter with auto-fix..."
	@golangci-lint run -v --fix --timeout $(LINT_TIMEOUT) ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy

deps-audit: ## Check dependencies for known vulnerabilities (govulncheck)
	@echo "Checking dependencies..."
	@govulncheck ./...
