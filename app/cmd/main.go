package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/qurk0/pr-service/internal/api/handlers"
	"github.com/qurk0/pr-service/internal/config"
	"github.com/qurk0/pr-service/internal/domain/services"
	"github.com/qurk0/pr-service/internal/metrics"
	"github.com/qurk0/pr-service/internal/storage/pgsql"
	"github.com/qurk0/pr-service/pkg/middlewares"
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
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to read configs: %v", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.LogLevel)
	mainLogger := logger.With("op", "cmd.main")
	mainLogger.Debug("debug messages are enable")

	mainLogger.Debug("creating storage instanse")
	db, err := pgsql.NewDB(context.Background(), cfg.ConnString())

	storage := pgsql.NewStorage(db, logger)

	servs := services.NewServices(storage.User, storage.Team, storage.PullRequest, logger)

	router := handlers.NewRouter(servs)

	metrics.Init()
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(middlewares.RequestLoggerMiddleware(logger))

	router.RegRoutes(app)

	if err := app.Listen(cfg.ListenAddr()); err != nil {
		logger.Error("failed to start listening addr", slog.String("err", err.Error()))
		os.Exit(1)
	}

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
