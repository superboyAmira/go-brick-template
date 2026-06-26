# Модули (`module/`)

**Модуль** — инфраструктура и обёртки внешних API. Не содержит бизнес-логики домена.

Доменная логика живёт в **кирпичах** → [BRICKS.md](./BRICKS.md).

---

## Интерфейс Module

Каждый модуль реализует `runtime.Module` из `internal/application/runtime`:

```go
type Module interface {
    Init(ctx context.Context, info Info) error
    Run(ctx context.Context, c *Closer)
}
```

Зависимости передаются через **конструктор** `New(...)`, не через глобальный реестр.

### Паттерн реализации

```go
func New(ctx context.Context, opts *options.PostgresOptions) (*Module, error) {
    db, err := openPools(ctx, opts)
    if err != nil {
        return nil, err
    }
    return &Module{db: db}, nil
}

func (m *Module) Init(ctx context.Context, _ runtime.Info) error {
    return nil // или поздняя инициализация без config lookup
}

func (m *Module) Run(_ context.Context, c *runtime.Closer) {
    c.Add("postgres", func() error { /* close pools */ return nil })
}
```

---

## Реализованные в шаблоне

| Имя | Путь | Конструктор | Назначение |
|-----|------|-------------|------------|
| `postgres` | `module/postgres/` | `New(ctx, *PostgresOptions)` | pgx pools |
| `http` | `module/http/` | `New(*HTTPOptions)` + `MountRoutes(cfg, *Registry)` | Public API |
| `admin_http` | `module/admin_http/` | `New(*HTTPOptions, *postgres.DB)` | Admin endpoints |

---

## Расширение (паттерн networking-backend)

1. `module/<name>/module.go` — `New(...)` с явными зависимостями, `Init`/`Run`
2. Регистрация в `internal/application/application.go` (`Run`)
3. Секция в `internal/config/options/` и `config.LoadFromEnv()`
4. Строка в `docs/MODULES.md` и `llm/manifest.json`

Типичные модули: `redis`, `qdrant`, `storage`, `llm`, `worker`, `apple`.

**Init order** (networking-backend, при расширении):

```
postgres → redis → qdrant → storage → llm → worker → http → admin_http
```

---

## Правила

1. Кирпичи **не** создают pgx/redis-клиенты напрямую
2. AI inference **только** через `module/llm`, `module/embedder`, `module/whisper`
3. Модули получают зависимости в `New()` или из `wire.go`, не из глобального state
4. `module/*` импортирует `internal/application/runtime` (интерфейс), не полный `application` (избегает import cycle)

---

## DB (postgres)

```go
func (db *DB) Write(ctx context.Context) *pgxpool.Pool { return db.master }
func (db *DB) Read(ctx context.Context) *pgxpool.Pool  { return db.slave }
```

Кирпичи используют DB только внутри `repository/postgres.go`.
