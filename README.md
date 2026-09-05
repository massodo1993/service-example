# Микросевисы

## Сервисы

| Сервис      | Транспорт | Описание                                  |
|-------------|-----------|-------------------------------------------|
| `order`     | HTTP (ogen)| Заказы: создание, оплата, отмена, чтение |
| `inventory` | gRPC       | Каталог деталей                          |
| `payment`   | gRPC       | Проведение платежей                     |

Слои в каждом сервисе: `internal/api` → `internal/service` → `internal/repository` (+ `internal/client` для внешних вызовов).

## Тесты

```bash
task test           # запустить unit-тесты всех модулей
task test-coverage  # тесты + процент покрытия по каждому модулю
task coverage:html  # HTML-отчёт покрытия
task mockery:gen    # перегенерировать моки (.mockery.yaml)
```

### Покрытие (`task test-coverage`)

| Модуль      | Покрытие |
|-------------|----------|
| `order`     | 81.2%    |
| `inventory` | 66.3%    |
| `payment`   | 100.0%   |
| **Итого**   | **74.9%**|

Покрываются ключевые usecase сервисного слоя (моки repo/client через `mockery`) и
in-memory репозитории.

## API-тесты

```bash
task test-api  # end-to-end проверка через gRPC + HTTP (сервисы должны быть запущены)
```
