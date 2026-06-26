# Кирпичи (`internal/brick/`)

**Кирпич** — изолированный домен с model, repository, service и публичным `contract.go`.

Инфраструктура — в **модулях** → [MODULES.md](./MODULES.md).

---

## Структура кирпича

```
internal/brick/item/
├── contract.go       # Service interface — единственная точка входа для других кирпичей
├── model/
│   └── item.go
├── repository/
│   └── postgres.go   # Squirrel + pgx через module/postgres DB
└── service/
    └── service.go    # implements item.Service (contract)
```

### contract.go

```go
package item

type Service interface {
    List(ctx context.Context) ([]Item, error)
    Get(ctx context.Context, id uuid.UUID) (*Item, error)
    Create(ctx context.Context, req CreateRequest) (*Item, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

YAML-дубликат: `docs/contracts/bricks/item.yaml`.

---

## Persistence (Squirrel)

```go
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
    q := psql.Select("id", "title", "created_at", "updated_at").
        From("items").
        Where(sq.Eq{"id": id})
    sql, args, err := q.ToSql()
    // row := r.db.Read(ctx).QueryRow(ctx, sql, args...)
}
```

| Операция | Pool |
|----------|------|
| SELECT | `db.Read(ctx)` |
| INSERT, UPDATE, DELETE | `db.Write(ctx)` |

- Схема: только goose `migrations/`
- **Не использовать:** sqlc, `gen/`, `query_*.sql`

---

## Правила изоляции

1. **Запрещён** прямой import другого кирпича
2. **Разрешены** только contract-интерфейсы (`item.Service`, …)
3. HTTP handlers вызывают **только** contract, не `service`/`repository` напрямую
4. Wiring — `internal/application/wire.go` + `internal/registry/registry.go`

---

## Пример кирпича в шаблоне

| Кирпич | Таблицы | Contract YAML |
|--------|---------|---------------|
| **item** | `items` | [bricks/item.yaml](./contracts/bricks/item.yaml) |

---

## Добавление нового кирпича

1. `internal/brick/<name>/` — contract, model, repository, service
2. `docs/contracts/bricks/<name>.yaml`
3. Пути в `docs/contracts/openapi.yaml` (contract-first)
4. `wire.go` — создать repo + service, положить в `bricks.Registry`
5. `internal/api/handlers/` — routes через contract
6. `migrations/` — owned tables

---

## Wiring (composition root)

```go
// internal/application/wire.go

itemRepo := itemrepo.NewPostgres(pg)
itemSvc := itemsvc.New(itemRepo)

return &registry.Registry{Item: itemSvc}
```

При split: `itemSvc` → HTTP/gRPC client, реализующий `item.Service` из YAML.

---

## Ссылки

- [contracts/bricks/_types.yaml](./contracts/bricks/_types.yaml) — shared types
- [ARCHITECTURE.md](./ARCHITECTURE.md)
