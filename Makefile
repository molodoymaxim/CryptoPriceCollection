# Переменные
APP_NAME=crypto
DOCKER_COMPOSE=docker-compose.yml
GO=go
GOFLAGS=-v
PROJECT_NAME=cryptopricecollection
VOLUME_NAME=$(PROJECT_NAME)_postgres-data

.PHONY: all build run test clean docker-build docker-up docker-down docker-clean migrate-up migrate-down

# Сборка и тестирование
all: docker-down docker-build docker-up

# Сборка бинарника Go
build:
	$(GO) build $(GOFLAGS) -o $(APP_NAME) ./cmd/main.go

# Локальный запуск приложения
run: build
	./$(APP_NAME)

# Очистка бинарников
clean:
	rm -f $(APP_NAME)

# Сборка Docker-образа
docker-build:
	docker-compose -f $(DOCKER_COMPOSE) build

# Запуск контейнеров в фоновом режиме
docker-up:
	docker-compose -f $(DOCKER_COMPOSE) up -d

# Остановка и удаление контейнеров
docker-down:
	docker-compose -f $(DOCKER_COMPOSE) down

# Очистка контейнеров и томов
docker-clean: docker-down
	docker volume rm $(VOLUME_NAME) || true

# Применение миграций (через запуск приложения)
migrate-up: docker-up
	@echo "Миграции применяются автоматически при запуске приложения"

# Откат последней миграции (локально)
migrate-down:
	$(GO) run -tags migrate ./cmd/main.go down

# Генерация Swagger-документации
swagger:
	swag init -g cmd/main.go -o docs