# Error contract

Все REST endpoints возвращают ошибки в едином формате.

## Envelope

```json
{
  "error": {
    "code": "ITEM_NOT_FOUND",
    "message": "Item not found",
    "details": {},
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `code` | string | Стабильный машиночитаемый код (SCREAMING_SNAKE_CASE) |
| `message` | string | Человекочитаемое описание |
| `details` | object | Доп. контекст (поля валидации, entity_id, …) |
| `request_id` | string | UUID из middleware `X-Request-ID` |

## HTTP status mapping

| Status | Когда |
|--------|-------|
| 400 | Невалидный запрос, validation |
| 404 | Сущность не найдена |
| 500 | Внутренняя ошибка |
| 503 | Модуль/внешний API недоступен |

## Коды (пример)

| Code | HTTP | Описание |
|------|------|----------|
| `VALIDATION_ERROR` | 400 | Невалидные поля |
| `ITEM_NOT_FOUND` | 404 | Item не найден |
| `INTERNAL_ERROR` | 500 | Неожиданная ошибка |
