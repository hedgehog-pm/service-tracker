# service-tracker

Ядро трекера задач на Go и PostgreSQL. Реализован CRUD проектов,
статусов, модулей и задач.

- [Проект решения](docs/design.md).
- [OpenAPI: все операции](api/server/openapi.yaml).
- Контракты ресурсов: [проекты](api/server/projects/openapi.yaml),
  [статусы](api/server/statuses/openapi.yaml),
  [модули](api/server/modules/openapi.yaml),
  [задачи](api/server/tasks/openapi.yaml).
- [Эталонная схема БД](docs/schema.sql).

Каталог создан через `gonew` из закреплённой версии болерплейта:

```sh
go run golang.org/x/tools/cmd/gonew@v0.52.0 \
  github.com/hedgehog-pm/service-bolerplate@v0.0.0-20261022140028-5b26cee97865 \
  github.com/hedgehog-pm/service-tracker ./service-tracker
```

Запуск из корня сервиса:

```sh
docker compose up -d --wait postgres
export DATABASE_URL='postgres://tracker:tracker@localhost:5433/tracker?sslmode=disable'
go run ./cmd/migrate
go run ./cmd/service
```

HTTP слушает `:8080`; адрес задаётся через `HTTP_ADDR`.
Миграции создают все четыре таблицы. Для каждого ресурса доступны
GET/POST коллекции и GET/PUT/DELETE объекта:

| Ресурс | Коллекция | Объект |
| --- | --- | --- |
| Проекты | `/projects` | `/projects/{project_id}` |
| Статусы | `/projects/{project_id}/statuses` | `/projects/{project_id}/statuses/{status_id}` |
| Модули | `/projects/{project_id}/modules` | `/projects/{project_id}/modules/{module_id}` |
| Задачи | `/projects/{project_id}/tasks` | `/projects/{project_id}/tasks/{task_id}` |

Все списки принимают limit и offset; список задач — ещё status_id и module_id.
Defaults limit=50 и offset=0 заданы в OpenAPI. Middleware подставляет их
до биндинга параметров, поэтому обработчики получают заполненные значения.
PUT задачи требует title, description, status_id и module_id; null снимает модуль.
Ключ проекта неизменяем. Удаление проекта каскадно удаляет его данные.
Используемые задачами статусы и модули удалить нельзя: API возвращает 409.

Пример для пустой базы; в следующих запросах используются ID из ответов:

```sh
curl -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"key":"WEB","name":"Web"}'

curl -X POST http://localhost:8080/projects/1/statuses \
  -H 'Content-Type: application/json' \
  -d '{"name":"Todo","category":"todo"}'

curl -X POST http://localhost:8080/projects/1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Первая задача","status_id":1}'
```

Генерация из корня сервиса:

```sh
make generate
```

OpenAPI генерируется для common, projects, statuses, modules и tasks
в internal/generated/api. SQL-запросы лежат в queries.sql рядом с каждым
репозиторием: project, status, module, task. Bob генерирует отдельный пакет
для каждого из них в internal/generated/bob.
Генераторы закреплены в go.mod и запускаются через go tool.

Для SQL-кодгена нужна БД с применёнными миграциями. По умолчанию Makefile
использует PostgreSQL из compose.yaml на порту 5433.
Адрес можно изменить через DATABASE_URL или отдельно через PSQL_DSN.
В PostgreSQL-запросах используются параметры $1, $2 и аннотации Bob `-- QueryName`.
Опциональные фильтры и пагинация добавляются билдером к сгенерированному ListTasks.
Результаты читаются сгенерированными mapper-ами Bob через pgxscan.One/All
и преобразуются явными конвертерами в доменные сущности.
Отдельной копии Task для чтения нет. Транзакцией управляет use case
через отдельный Transactor из internal/storage/postgres.

Числовые поля ядра в PostgreSQL используют INT, в Go — int.
В OpenAPI задан format: int32 и x-go-type: int; Bob получает int
через алиасы в bobgen.yaml. Limit и offset имеют одинаковый тип.

Общий перевод ошибок PostgreSQL находится в internal/storage/postgres/errors.go.
Известные SQLSTATE превращаются в доменные ошибки, остальные —
в ErrInternalServerError. Имена ограничений не используются; техническая
причина остаётся только в логе.

Проверки:

```sh
make fmt lint
go build ./cmd/...
```
