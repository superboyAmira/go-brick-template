# Go Brick Template

Шаблон Go-монолита с архитектурой **модули** (`module/`) + **кирпичи** (`internal/brick/`), скопированной из [networking-backend](https://github.com/networking/networking-backend).

**Go:** 1.23+  
**Persistence:** PostgreSQL + [goose](https://github.com/pressly/goose) + [Squirrel](https://github.com/Masterminds/squirrel) + pgx

## Документация

| Документ | Описание |
|----------|----------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Слои, wiring, DI |
| [docs/MODULES.md](docs/MODULES.md) | Инфра (`module/`) |
| [docs/BRICKS.md](docs/BRICKS.md) | Доменные кирпичи (`internal/brick/`) |
| [docs/contracts/README.md](docs/contracts/README.md) | Индекс контрактов |
| [docs/contracts/openapi.yaml](docs/contracts/openapi.yaml) | REST API v1 |
| [llm/manifest.json](llm/manifest.json) | Карта проекта для LLM-агентов |
| [llm/codebase/manifest.json](llm/codebase/manifest.json) | Карта кодовой базы |

## Cursor

| Ресурс | Описание |
|--------|----------|
| [.cursor/rules/llm-manifest-navigation.mdc](.cursor/rules/llm-manifest-navigation.mdc) | Навигация по manifest для агентов |
| [.cursor/skills/go-brick-architecture/SKILL.md](.cursor/skills/go-brick-architecture/SKILL.md) | Skill: modules + bricks |

## Терминология

- **Модуль** — инфраструктура и обёртки внешних API (`module/postgres`, `module/http`, …)
- **Кирпич** — изолированный домен (`internal/brick/item`, …), связь только через `contract.go`

## Структура

```
go-brick-template/
├── cmd/app/              # entrypoint
├── module/               # модули (infra)
├── internal/
│   ├── brick/            # кирпичи (домен)
│   ├── api/              # HTTP handlers
│   ├── application/      # composition root (Run, wire, runtime)
│   ├── registry/         # wired brick contracts
│   ├── config/
│   ├── shared/              # apperr, shared enums
│   └── pkg/                 # generic utils (pagination)
├── migrations/           # goose
├── docs/contracts/       # OpenAPI + inter-brick YAML
└── llm/                  # manifests для агентов
```

## Быстрый старт

```bash
make dev                    # postgres + goose migrate
cp deploy/.env.example .env # загрузить в IDE / shell
make run                    # API :8080, admin :9090
```

**Swagger UI:**

| URL | Описание |
|-----|----------|
| http://localhost:8080/swagger | REST API `/api/v1` |
| http://localhost:8080/openapi.yaml | Сырой spec |
| http://localhost:9090/swagger | Admin (`/health`, `/ready`) |

**Пример API:**

```bash
curl -s localhost:8080/api/v1/items | jq
curl -s -X POST localhost:8080/api/v1/items -d '{"title":"hello"}' | jq
```

## Contract-first

Правки REST только в `docs/contracts/openapi.yaml`, затем:

```bash
make swagger-gen
make swagger-check   # CI
make build
```

## Makefile

```bash
make help
make validate-spec
make migrate-up
make migrate-create NAME=add_feature
make test
make lint
```

Переменные окружения:

| Переменная | По умолчанию |
|------------|--------------|
| `DATABASE_URL` | `postgres://brick:brick@localhost:5432/brick?sslmode=disable` |
| `HTTP_ADDR` | `:8080` |
| `ADMIN_HTTP_ADDR` | `:9090` |

## Добавление кирпича

См. [docs/BRICKS.md](docs/BRICKS.md) и skill [.cursor/skills/go-brick-architecture/SKILL.md](.cursor/skills/go-brick-architecture/SKILL.md).

1. `internal/brick/<name>/` — contract, model, repository, service
2. `docs/contracts/bricks/<name>.yaml`
3. `wire.go` + `application.Registry`
4. handlers + `openapi.yaml`
5. `migrations/`

## License

See [LICENSE](LICENSE).
