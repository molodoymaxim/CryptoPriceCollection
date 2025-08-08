# CryptoPriceCollection
## Описание проекта
CryptoPriceCollection — это REST API на Go для отслеживания цен криптовалют с использованием внешнего API CoinGecko. Приложение позволяет:
- Добавлять криптовалюты в список отслеживаемых
- Удалять криптовалюты из списка, сохраняя их исторические цены в базе данных
- Получать последнюю цену валюты или цену, ближайшую к указанному времени
- Хранить данные в PostgreSQL с автоматическими миграциями
- Просматривать документацию API через Swagger UI

## Функциональность
- **Эндпоинты API**:
  - `POST /currency/add` — добавляет валюту в отслеживаемый список
  - `POST /currency/remove` — удаляет валюту из списка
  - `POST /currency/price` — возвращает последнюю цену (без `timestamp`) или ближайшую цену к указанному времени (с `timestamp`)
- **Swagger UI**: Документация API доступна по адресу `http://host:port/swagger/index.html`
- **Фоновый процесс**: Получение цен от CoinGecko API каждые N секунд (настраивается через `FETCH_INTERVAL`)
- **База данных**: PostgreSQL с таблицами `watched_currencies` и `currency_prices`

## Установка и запуск
### 1. Клонирование репозитория
```bash
git clone <repository_url>
cd CryptoPriceCollection
```

### 2. Настройка .env
Создайте файл `.env` в корне проекта на основе `.env.example`:
```plaintext
# Конфигурация PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=db_name
POSTGRES_SSLMODE=disable
POSTGRES_QUERY_TIMEOUT=5

# Конфигурация пула соединений базы данных
DB_MAX_CONN=4
DB_CONN_IDLE_TIME=600
DB_CONN_LIFE_TIME=1800

# Конфигурация HTTP-сервера
HTTP_PORT=1234
HTTP_READ_TIMEOUT=10
HTTP_WRITE_TIMEOUT=10
HTTP_IDLE_TIMEOUT=60
HTTP_SHUTDOWN_TIMEOUT=5

# Конфигурация клиента внешнего API (CoinGecko)
API_BASE_URL=https://api.coingecko.com/api/v3
API_TIMEOUT=10
API_MAX_RETRIES=3
API_RETRY_BACKOFF=500

# Конфигурация фоновых задач
FETCH_INTERVAL=60000
BATCH_INTERVAL=60000
```

### 3. Установка зависимостей
```bash
go mod tidy
```

### 4. Генерация Swagger-документации
```bash
make swagger
```

### 5. Запуск приложения
```bash
make all
```
Также есть возможность использовать другие команды из Makefile для разных этапов запуска приложения

### 6. Доступ к API
- API доступно по адресу: `http://host:port`
- Swagger UI: `http://host:port/swagger/index.html`

### 7. Тестирование
Примеры тестовых запросов:
- `POST /currency/add` с `{"coin": "bitcoin"}`
- `POST /currency/price` с `{"coin": "bitcoin"}`
- `POST /currency/price` с `{"coin": "bitcoin", "timestamp": 1754645360}`
- `POST /currency/remove` с `{"coin": "bitcoin"}`

### 8. Остановка приложения
```bash
make docker-down
make docker-clean
```

## Структура проекта
- `cmd/main.go` — точка входа приложения.
- `internal/handlers/` — обработчики HTTP-запросов.
- `internal/services/` — бизнес-логика.
- `internal/postgresql/` — работа с PostgreSQL.
- `internal/types/` — структуры данных.
- `internal/server/` — настройка HTTP-сервера.
- `pkg/migrations/` — SQL-миграции для базы данных.
- `docs/` — сгенерированная Swagger-документация.

## Ограничения
- **CoinGecko API**: Бесплатная версия имеет лимит 30–50 запросов/мин. При превышении возвращается ошибка 429. Настройте `FETCH_INTERVAL` ≥ 60000 для минимизации проблем.
