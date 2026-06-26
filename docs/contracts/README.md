# Contracts index

Спецификации контрактов go-brick-template.

## Терминология

| Термин | Документ | Назначение |
|--------|----------|------------|
| **Модуль** | [MODULES.md](../MODULES.md) | Инфра + внешние API |
| **Кирпич** | [BRICKS.md](../BRICKS.md) | Доменная логика |

## Внешние контракты

| Файл | Описание |
|------|----------|
| [openapi.yaml](./openapi.yaml) | REST API v1 (OpenAPI 3.1) |
| [errors.md](./errors.md) | Формат ошибок, коды |

**Base URL:** `/api/v1`  
**Extensions:** `x-brick`, `x-modules`

**Swagger UI** (при запущенном `go run ./cmd/app`):

- http://localhost:8080/swagger
- http://localhost:9090/swagger — admin endpoints

**Автогенерация Swagger** из контракта (источник правды):

```bash
make swagger-gen
make build          # зависит от swagger-gen
make swagger-check  # CI: embed совпадает с openapi.yaml
```

## Inter-brick contracts

| Кирпич | YAML |
|--------|------|
| item (example) | [bricks/item.yaml](./bricks/item.yaml) |

Shared types: [bricks/_types.yaml](./bricks/_types.yaml)

## Database

| Файл | Описание |
|------|----------|
| [../../migrations/20260626120000_init_items.sql](../../migrations/20260626120000_init_items.sql) | Goose init (example `items` table) |

## LLM agents

См. [../../llm/manifest.json](../../llm/manifest.json) и [../../llm/codebase/manifest.json](../../llm/codebase/manifest.json).
