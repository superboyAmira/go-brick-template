---
name: go-brick-architecture
description: >-
  Go monolith architecture with modules (module/) and domain bricks
  (internal/brick/). Use when adding bricks, modules, HTTP handlers, migrations,
  or wiring in go-brick-template / networking-backend style projects.
---

# Go Brick Architecture (modules + bricks)

Архитектура из networking-backend: монолит готов к split в микросервисы.

## When to Activate

- Добавление доменной логики (кирпич)
- Новый REST endpoint
- Инфра-модуль (postgres, redis, llm, worker)
- Миграции и repository
- Wiring в application layer

## Терминология

| Термин | Путь | Роль |
|--------|------|------|
| **Модуль** | `module/<name>/` | Инфра и внешние API; `app.Module` |
| **Кирпич** | `internal/brick/<name>/` | Домен: model + repository + service + `contract.go` |

**Модуль ≠ кирпич.**

## Слои (сверху вниз)

1. `module/http` → `internal/api/handlers` (Fiber)
2. Handlers вызывают **только** `contract.Service`
3. Service → repository → `module/postgres.DB`
4. Cross-brick: только contract interfaces в `wire.go`

## Структура кирпича

```
internal/brick/<name>/
├── contract.go
├── model/
├── repository/postgres.go   # Squirrel + pgx
└── service/service.go
```

## Wiring checklist

1. Repository: `NewPostgres(pg *postgres.DB)`
2. Service: `New(repo, ...contracts)`
3. `internal/application/wire.go` → `registry.Registry`
4. Handler в `internal/api/handlers/`
5. `docs/contracts/openapi.yaml` (contract-first)
6. `make swagger-gen`
7. Goose migration для owned tables

## Persistence rules

- Схема: **goose** в `migrations/`
- Запросы: **Squirrel** + pgx
- SELECT → `db.Read(ctx)`, INSERT/UPDATE/DELETE → `db.Write(ctx)`
- **Не использовать sqlc**

## Isolation rules

- ❌ `import ".../brick/other"` из другого кирпича
- ✅ `other.Service` interface из contract
- ❌ pgx/redis клиенты вне `module/`
- ❌ global `appstate` — DI через конструкторы в `application/`
- ❌ handlers → repository напрямую

## Task routing

| Задача | Файлы |
|--------|-------|
| REST route | `handlers/`, `openapi.yaml`, `swagger-gen` |
| Domain logic | `brick/<name>/service/`, `contract.go` |
| SQL table | `migrations/`, `repository/` |
| External API | `module/<name>/module.go`, `config/options` |

## Navigation

Перед правками читай:

1. `llm/codebase/manifest.json`
2. `docs/ARCHITECTURE.md`
3. Точечно файлы из `task_routing`

## Reference

Полный production-пример: networking-backend (`internal/brick/*`, `module/worker`, async jobs).
