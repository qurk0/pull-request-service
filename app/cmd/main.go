package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/config"
	"github.com/qurk0/pr-service/internal/domain/services"
	"github.com/qurk0/pr-service/internal/storage/pgsql"
)

const (
	// Путь до файла конфигурации
	cfgPath = "configs/cfg.yaml"

	// Уровни логирования
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func main() {
	// TODO: Читаем конфиги (либа cleanenv)
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to read configs: %w", err)
		os.Exit(1)
	}

	// Запускаем логгер
	// В рамках тестового задания это будет log/slog с выводом в терминал
	logger := newLogger(cfg.LogLevel)
	mainLogger := logger.With("op", "cmd.main")
	mainLogger.Debug("debug messages are enable")

	// Создаём инстанс стореджа
	mainLogger.Debug("creating storage instanse")
	db, err := pgsql.NewDB(context.Background(), cfg.ConnString())

	// Тут создаём структуру Storage, которая хранит в себе 3 репозитория для работы с юзерами, тимами и ПР'ами, всё в одном месте для удобства
	storage := pgsql.NewStorage(db)

	// Создаём инстанс сервисов, storage будет реализовывать методы интерфейсов сервисов
	servs := services.NewServices(storage)

	// TODO: Создаём fiber.App и привязываем хэндлеры к эндпоинтам
	app := fiber.New()

	// TODO: Слушаем адрес из конфигов

	// TODO: Graceful Shutdown
}

// По умолчанию уровень логирования - info
// Подробнее про уровня логирования - в README.md
func newLogger(logLevel string) *slog.Logger {
	switch logLevel {
	case LevelDebug:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case LevelWarn:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))

	case LevelError:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))

	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
}
