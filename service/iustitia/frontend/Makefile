.DEFAULT_GOAL := help
SHELL         := bash
.SHELLFLAGS   := -c

# ------------------------------------------------------------------------------
# Paths
# ------------------------------------------------------------------------------
OPENAPI_SRC     := ../backend/internal/openapi/openapi.yml
API_SCHEMA      := src/shared/api/schema.d.ts
DIST_DIR        := dist
IMAGE_NAME      ?= iustitia-iustitia-frontend
IMAGE_TAG       ?= latest
COMPOSE_FILE    := ../docker-compose.yml
COMPOSE_SERVICE := iustitia-frontend

# ------------------------------------------------------------------------------
# HELP
# ------------------------------------------------------------------------------
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup

.PHONY: install

install: ## Install npm dependencies from package-lock.json
	@echo "Installing dependencies..."
	@npm ci
	@echo "Dependencies installed"

##@ Code generation

.PHONY: gen-api

gen-api: ## Generate TypeScript types from backend OpenAPI spec
	@echo "Generating API schema from $(OPENAPI_SRC)..."
	@npm run gen:api
	@echo "Schema generated: $(API_SCHEMA)"

##@ Build & Run

.PHONY: dev build preview clean

dev: ## Run Vite dev server (localhost:3000)
	@echo "Starting dev server..."
	@npm run dev

build: ## Build production bundle (typecheck + vite build)
	@echo "Building production bundle..."
	@npm run build
	@echo "Built: $(DIST_DIR)"

preview: ## Preview production build locally
	@npm run preview

clean: ## Remove build artefacts and node_modules cache
	@echo "Cleaning..."
	@rm -rf $(DIST_DIR) node_modules/.vite .tsbuildinfo
	@echo "Clean completed"

##@ Tests

.PHONY: test

test: ## Run vitest suite
	@echo "Running tests..."
	@npm run test

##@ Quality

.PHONY: fmt lint lint-fix typecheck

fmt: ## Format code with prettier
	@echo "Formatting..."
	@npm run format
	@echo "Done"

lint: ## Run eslint
	@echo "Running linter..."
	@npm run lint

lint-fix: ## Run eslint with --fix
	@echo "Running linter with auto-fix..."
	@npm run lint:fix

typecheck: ## Run tsc --noEmit
	@echo "Typechecking..."
	@npm run typecheck

##@ Docker

.PHONY: docker-build docker-up docker-down docker-logs docker-restart

docker-build: ## Build frontend docker image via compose (picks up build.args)
	@echo "Building docker image..."
	@docker compose -f $(COMPOSE_FILE) build $(COMPOSE_SERVICE)

docker-up: ## Start frontend container (detached)
	@docker compose -f $(COMPOSE_FILE) up -d $(COMPOSE_SERVICE)

docker-down: ## Stop and remove frontend container
	@docker compose -f $(COMPOSE_FILE) rm -sf $(COMPOSE_SERVICE)

docker-logs: ## Tail frontend container logs
	@docker compose -f $(COMPOSE_FILE) logs -f $(COMPOSE_SERVICE)

docker-restart: docker-down docker-up ## Restart frontend container
