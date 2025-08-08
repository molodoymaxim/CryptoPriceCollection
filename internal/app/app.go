package app

import (
	"CryptoPriceCollection/internal/config"
	"CryptoPriceCollection/internal/handlers"
	"CryptoPriceCollection/internal/logger"
	"CryptoPriceCollection/internal/repositories"
	"CryptoPriceCollection/internal/server"
	"CryptoPriceCollection/internal/services"
	"CryptoPriceCollection/internal/system"
	"CryptoPriceCollection/internal/types"
	"context"
	"fmt"
	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"time"
)

func Start() {
	// Создаем логгер
	logCust := logger.New()
	logCust.SetServiceName("CryptoPriceCollection")

	// Конфигурации
	cfgPostgres := &types.ConfigPostgres{}
	cfgConnDB := &types.ConfigConnDB{}
	cfgHTTPServer := &types.ConfigHTTPServer{}
	cfgAPIClient := &types.ConfigAPIClient{}
	cfgTasks := &types.ConfigTasks{}

	// Подгружаем конфигурацию из переменных окружения
	err := config.GetConfigsPath([]any{
		cfgPostgres,
		cfgConnDB,
		cfgHTTPServer,
		cfgAPIClient,
		cfgTasks,
	})
	if err != nil {
		logCust.WriteLog(logrus.FatalLevel, "Get config in enviroment var", logrus.Fields{
			"func":       "config.GetConfigsPath",
			"error":      err,
			"stacktrace": fmt.Sprintf("%+v", errors.WithStack(err)),
		})
	}

	// Формируем общий конфиг всего микросервиса
	cfgApp := &types.ConfigApp{
		Postgres:   *cfgPostgres,
		ConnDB:     *cfgConnDB,
		HTTPServer: *cfgHTTPServer,
		APIClient:  *cfgAPIClient,
		Tasks:      *cfgTasks,
	}

	// Устанавливаем формат логов как GELF
	gelfFmt := formatter.NewGelf("CryptoPriceCollection")
	logCust.SetFormater(gelfFmt)

	// Формирование строки подключения для миграций
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(cfgApp.Postgres.PostgresUser),
		url.QueryEscape(cfgApp.Postgres.PostgresPassword),
		cfgApp.Postgres.PostgresHost,
		cfgApp.Postgres.PostgresPort,
		cfgApp.Postgres.PostgresDBName,
		cfgApp.Postgres.SSLMode)

	// Инициализация миграций
	m, err := migrate.New("file://pkg/migrations/", connStr)
	if err != nil {
		log.Fatalf("Migration initialization error: %v", err)
	}
	defer m.Close()

	// Проверка аргументов командной строки для миграций
	if len(os.Args) > 1 && os.Args[1] == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration Rollback error: %v", err)
		}
		log.Println("Migrations successfully rolled out")
		return
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration execution error: %v", err)
	}

	// Инициализация зависимостей системы
	syst, err := system.New(&cfgApp.Postgres, &cfgApp.ConnDB)
	if err != nil {
		logCust.WriteLog(logrus.FatalLevel, "Create system", logrus.Fields{
			"func":       "system.New",
			"error":      err,
			"stacktrace": fmt.Sprintf("%+v", errors.WithStack(err)),
		})
	}
	logCust.WriteLog(logrus.InfoLevel, "Successful create system", logrus.Fields{})

	// Инициализация репозитория
	repo := repositories.New(syst)

	// Инициализация сервиса
	service := services.NewService(*repo, cfgAPIClient.BaseURL, time.Duration(cfgTasks.FetchInterval), time.Duration(cfgTasks.BatchInterval))

	// Выборка цен в фоновом режиме
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.CryptoService.StartPriceFetcher(ctx)

	// Инициализация ручек
	handler := handlers.NewHandler(service)

	// Инициализация роутера
	router := handler.InitRoutes()

	// Инициализация и запуск сервера
	s := server.NewServer(cfgApp, router)
	if err := s.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
