# {{cookiecutter.project_name}}

{{cookiecutter.project_description}}

### Установка зависимостей

Предварительные требования: Go 1.19 и выше

```bash
go env -w GOPRIVATE="git-ffd.kz/*"
go mod tidy
```

### Запуск

1. Создать в PostgreSQL схему `update_template` 
   (если используется не PostgreSQL - в [internal/app/app.go](internal/app/app.go) закомментить/удалить
    строки запуска миграций - `migrate.Up`)
2. Скопировать `example.env` и переименовать в `.env`
3. Заполнить в `.env` файле переменные `VAULT_USER` и `VAULT_PASSWORD`
4. Запустить проект 
```bash
go run main.go
```

- Документация:
  [http://localhost:{{cookiecutter.port}}/{{cookiecutter.project_name | slugify | lower}}-docs/docs/index.html](http://localhost:{{cookiecutter.port}}/{{cookiecutter.project_name | slugify | lower}}-docs/docs/index.html)
- Prometheus метрики:
  [http://localhost:{{cookiecutter.port}}/actuator/prometheus](http://localhost:{{cookiecutter.port}}/actuator/prometheus)
- Ping: [http://localhost:{{cookiecutter.port}}/ping](http://localhost:{{cookiecutter.port}}/ping)
- Healthcheck: [http://localhost:{{cookiecutter.port}}/readiness](http://localhost:{{cookiecutter.port}}/readiness)

### Дополнительные пакеты

В шаблоне используется ряд внутренних пакетов, о правильном использовании которых стоит прочитать
перед началом работы с шаблоном

* Логгер - [golog](https://git-ffd.kz/pkg/golog)
* Транзакции базы данных - [gotransaction](https://git-ffd.kz/pkg/gotransaction)
* Ошибки и их обработка - [ferr](https://git-ffd.kz/fmobile/ferr) для Fmobile
   * Основан на [goerr](https://git-ffd.kz/pkg/goerr)

### Конфигурация

Приложение конфигурируется при помощи переменных среды и `.env` файла.

#### Первоначальная настройка

1. Скопировать [`example.env`](example.env) файл и переименовать его в `.env`
2. Открыть `.env` файл исправить переменные, если это необходимо, заполнить `VAULT_USER` и `VAULT_PASSWORD`

#### Добавление переменных

1. Добавить переменную в [`example.env`](example.env) файл. В качестве
   значения установить оптимальное значение для локальной разработки/дев-среды,
   либо оставить пустым. При необходимости оставить комментарий с описанием
   переменной. Комментарии в `.env` файлах пишутся после знака `#`.
2. Добавить эту же переменную в `.env` файл
3. Добавить переменную в конфиг [internal/config/config.go](internal/config/config.go)
   по аналогии с уже имеющимися. При указании тега `envconfig` переменная спарсится
   в конфиг из `.env` файла автоматически.

### model и schema

В примере используется разделение на 2 пакета `model` и `schema`.

В `model` хранятся структуры, которые хранятся в базе данных.

В `schema` хранятся структуры для запросов/ответов которые используются в хендлерах.

В простых кейсах, когда модель базы данных полностью или почти полностью совпадает с тем, что
должно возвращаться и приниматься от фронта, абсолютно нормально использовать структуру из `models`.

Но в случаях, когда какие-то поля не должны возвращаться фронту или входные данные требуют специфичного
формата - лучше выделить структуру в `schema`.

### Миграции

В проекте нет автоматического средства управления миграциями. Все
миграции пишутся вручную. В качестве средства применения миграциями
используется [goose](https://github.com/pressly/goose).

**ВАЖНО**: Схему для проекта в базе данных нужно будет создать вручную!

При каждом запуске приложения применяются миграции, обновляя их до
последней доступной версии!

##### Установка

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

##### Создание новой миграции

Перейти в папку [migrations](migrations) и запустить там следующую команду.

```bash
goose create MIGRATION_NAME sql
```

`MIGRATION_NAME` - название миграции.

В папке `migrations` создается новый файл. Со следующим содержимым
```sql
-- +goose Up
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
```

В первом блоке `-- +goose Up` между `-- +goose StatementBegin` и `-- +goose StatementEnd`
пишется миграция для поднятия версии базы данных.

Во втором блоке `-- +goose Down` между `-- +goose StatementBegin` и `-- +goose StatementEnd`
пишется миграция для отмены этих изменений.

##### Ручное управление миграциями

Можно также управлять миграциями через `cli`.

```bash
# Применение миграций до последней доступной версии
goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} up

# Отмена миграций до первой версии
goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} down

# Поднять до определенной версии
goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} up-to ${VERSION}

# Отпустить до определенной версии
goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} down-to ${VERSION}
```

- `${DATABASE_SCHEMA}` - схема базы данных в которой работает приложение
- `${DATABASE_DSN}` - строка для подключения к базе данных
- `${VERSION}` - версия миграции (таймстемп в файле миграции)

### Sentry

В темплейте есть интеграция с Sentry. Активируется она когда в `.env` будет указан
`SENTRY_DSN` - это подключение к проекту в Sentry.

Чтобы его получить нужно создать проект под новый сервис в Sentry, выбрать Go в качестве языка,
ввести название и после этого скопировать получившийся DSN.

### Swagger документация

Swagger документация генерируется на основе специальных комментариев перед
хендлером. Для этого используется [swag](https://github.com/swaggo/swag).

##### Установка

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

##### Генерация документации

При изменениях в этих комментариях/добавлении/удалении хендлеров необходимо
перегенирировать документацию. Для этого в корне проекта

```bash
go generate ./...
```

Документация будет доступна по адресу:
[http://localhost:{{cookiecutter.port}}/{{cookiecutter.project_name | slugify | lower}}-docs/docs/index.html](http://localhost:{{cookiecutter.port}}/{{cookiecutter.project_name | slugify | lower}}-docs/docs/index.html)
