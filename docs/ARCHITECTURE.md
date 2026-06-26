# Архитектура Go Brick Template

**Версия:** 1.0 | Июнь 2026

Монолит Go 1.23+, спроектированный по архитектуре [networking-backend](https://github.com/networking/networking-backend): **модули** (`module/`) + **кирпичи** (`internal/brick/`). Готов к возможному разбиению на микросервисы без переписывания доменной логики.

---

## Терминология

| Термин | Папка | Ответственность |
|--------|-------|-----------------|
| **Модуль** | `module/<name>/` | Инфра (БД, Redis, Qdrant) и обёртки внешних API (LLM, Zoom, Apple…) |
| **Кирпич** | `internal/brick/<name>/` | Домен: model + repository + service + `contract.go` |

**Модуль ≠ кирпич.** Не смешивать в коде и документации.

---

## Слои

```
┌─────────────────────────────────────────────────────────┐
│  Clients (HTTP, mobile, CLI)                            │
└───────────────────────────┬─────────────────────────────┘
                            │ REST /api/v1
                            ▼
┌─────────────────────────────────────────────────────────┐
│  module/http  →  internal/api/handlers                  │
│  (Fiber routes, middleware)                             │
└───────────────────────────┬─────────────────────────────┘
                            │ вызывает contract кирпича
                            ▼
┌─────────────────────────────────────────────────────────┐
│  internal/brick/*                                       │
│  contract.go ← единственная точка входа между кирпичами  │
└───────────┬─────────────────────────────┬───────────────┘
            │                             │
            ▼                             ▼
┌───────────────────────┐     ┌───────────────────────────┐
│  module/postgres      │     │  module/llm, storage, …   │
│  (Squirrel + pgx)     │     │  (внешние API)            │
└───────────────────────┘     └───────────────────────────┘
```

---

## Composition root

**`internal/application/`** — DI через конструкторы, без глобального `appstate`.

| Файл | Роль |
|------|------|
| `application.go` | `Run()` — composition entry |
| `wire.go` | Repository + service кирпичей, `Registry` |
| `registry.go` | в `internal/registry/` — тип `Registry` (избегает import cycle с handlers) |
| `runtime/app.go` | `Module` interface, lifecycle (`Application`) |
| `runtime/closer.go` | Graceful shutdown |

Поток запуска:

1. `config.LoadFromEnv()`
2. `postgres.New(ctx, cfg.Postgres)` — пулы в конструкторе
3. `wireBricks(pg.DB())` → `*Registry`
4. `http.New(cfg)`, `adminhttp.New(cfg, pg)` — зависимости в конструкторах
5. `runtime.Build(ctx)` → `Init()` модулей
6. `http.MountRoutes(cfg, reg)` — handlers получают registry явно
7. `runtime.Run(ctx, closer)`

---

## Порядок Init модулей (шаблон)

```
postgres → http → admin_http
```

Расширенный порядок (как в networking-backend): postgres → redis → qdrant → storage → llm → worker → http → admin_http. См. [MODULES.md](./MODULES.md).

---

## Persistence

| Компонент | Технология |
|-----------|------------|
| Схема БД | goose (`migrations/`) |
| Запросы | Squirrel + pgx |
| Read/Write | `module/postgres.DB` |

**sqlc не используется.**

---

## Контракты

| Тип | Файл | Потребитель |
|-----|------|-------------|
| REST (внешний) | `docs/contracts/openapi.yaml` | клиенты |
| Inter-brick | `docs/contracts/bricks/*.yaml` | разработчики, будущие микросервисы |
| Ошибки | `docs/contracts/errors.md` | все handlers |

---

## Microservice-ready

При выделении кирпича в сервис:

1. `contract.go` → gRPC/HTTP client с тем же контрактом из `docs/contracts/bricks/<name>.yaml`
2. Owned tables переезжают с кирпичом
3. Межсервисные вызовы — только через contract, не shared DB

---

## Ссылки

- [MODULES.md](./MODULES.md)
- [BRICKS.md](./BRICKS.md)
- [contracts/README.md](./contracts/README.md)
