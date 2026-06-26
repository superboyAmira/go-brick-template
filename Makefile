SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

UNAME_S := $(shell uname -s)
GOPATH_BIN := $(shell go env GOPATH 2>/dev/null || echo "$(HOME)/go")/bin
ifeq ($(UNAME_S),Darwin)
HOMEBREW_BIN := $(shell brew --prefix 2>/dev/null)/bin
export PATH := $(GOPATH_BIN):$(HOMEBREW_BIN):$(PATH)
else
export PATH := $(GOPATH_BIN):$(PATH)
endif

# --- Project ---
APP_NAME        ?= app
CMD_DIR         := ./cmd/$(APP_NAME)
BIN_DIR         := ./bin
BINARY          := $(BIN_DIR)/$(APP_NAME)
MIGRATIONS_DIR  ?= migrations
GO              ?= go
GOFLAGS         ?=
LDFLAGS         ?=

# --- Tooling ---
GOOSE   ?= $(GO) run github.com/pressly/goose/v3/cmd/goose@v3.24.1

# --- OpenAPI / Swagger ---
OPENAPI_SRC   ?= docs/contracts/openapi.yaml
OPENAPI_EMBED ?= internal/api/swagger/openapi.yaml
ADMIN_SPEC    ?= internal/api/swagger/admin-openapi.yaml
SWAGGER_GEN   ?= ./scripts/swagger-gen.sh

# --- Database ---
DATABASE_URL    ?= postgres://brick:brick@localhost:5432/brick?sslmode=disable
GOOSE_DRIVER    ?= postgres

HAS_GO_SRC := $(shell find . -name '*.go' -not -path './vendor/*' 2>/dev/null | head -1)

.PHONY: help
help: ## Показать цели
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1mGo Brick Template\033[0m\n\n"} \
		/^[a-zA-Z0-9_.-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 } \
		END { printf "\n" }' $(MAKEFILE_LIST)

.PHONY: install-tools
install-tools: ## Установить dev-инструменты
	$(GO) install github.com/pressly/goose/v3/cmd/goose@latest
	$(GO) install go.uber.org/mock/mockgen@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest
ifeq ($(UNAME_S),Darwin)
	@command -v golangci-lint >/dev/null 2>&1 || brew install golangci-lint
else
	@command -v golangci-lint >/dev/null 2>&1 || ( \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(GOPATH_BIN) latest )
endif

.PHONY: deps
deps: ## go mod download / tidy
	$(GO) mod download
	$(GO) mod tidy

.PHONY: build
build: swagger-gen ## Собрать бинарник в bin/
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD_DIR)
	@echo "built $(BINARY)"

.PHONY: run
run: ## Запустить приложение (go run)
	$(GO) run $(GOFLAGS) $(CMD_DIR)

.PHONY: clean
clean: ## Удалить bin/ и кэш тестов
	@rm -rf $(BIN_DIR)
	@$(GO) clean -testcache 2>/dev/null || true

.PHONY: fmt
fmt: ## gofmt + goimports
	@if [ -z "$(HAS_GO_SRC)" ]; then echo "no Go files"; else \
		gofmt -w -s $$(find . -name '*.go' -not -path './vendor/*'); \
		goimports -w $$(find . -name '*.go' -not -path './vendor/*'); fi

.PHONY: lint
lint: ## golangci-lint
	@if [ -z "$(HAS_GO_SRC)" ]; then echo "skip lint"; else \
		golangci-lint run --config .golangci.yaml --fix ./...; fi

.PHONY: test
test: ## go test ./...
	$(GO) test $(GOFLAGS) -race -count=1 ./...

.PHONY: check
check: swagger-check fmt lint test ## swagger + fmt + lint + test

.PHONY: migrate-up
migrate-up: ## goose up
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## goose down (1 шаг)
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" down

.PHONY: migrate-status
migrate-status: ## goose status
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" status

.PHONY: migrate-create
migrate-create: ## goose create (NAME=description)
	@test -n "$(NAME)" || (echo 'usage: make migrate-create NAME=add_feature'; exit 1)
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(NAME) sql

.PHONY: migrate-validate
migrate-validate: ## goose validate
	$(GOOSE) -dir $(MIGRATIONS_DIR) validate

.PHONY: swagger-gen
swagger-gen: ## Валидация контракта + embed openapi для Swagger UI
	@bash $(SWAGGER_GEN)

.PHONY: swagger-check
swagger-check: ## CI: контракт валиден и embed совпадает
	@bash $(SWAGGER_GEN)
	@diff -q $(OPENAPI_SRC) $(OPENAPI_EMBED) >/dev/null || \
		(echo "embed out of sync — run make swagger-gen"; exit 1)
	@echo "swagger-check: OK"

.PHONY: validate-spec
validate-spec: migrate-validate swagger-gen ## Миграции + OpenAPI

.PHONY: docker-up
docker-up: ## docker compose local infra up -d
	docker compose -f docker-compose.yml -f docker-compose.local.yml up -d

.PHONY: docker-down
docker-down: ## docker compose local down
	docker compose -f docker-compose.yml -f docker-compose.local.yml down

.PHONY: docker-migrate
docker-migrate: ## one-shot goose migrate via compose
	docker compose -f docker-compose.yml -f docker-compose.local.yml run --rm migrate

.PHONY: dev
dev: docker-up ## Local postgres + migrate
	@echo "Waiting for migrate..."
	@docker compose -f docker-compose.yml -f docker-compose.local.yml run --rm migrate
	@echo "Infra ready. Run: make run"
