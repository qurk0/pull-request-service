package middlewares

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLoggerMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		reqIDraw := c.Locals("requestid")
		reqID, _ := reqIDraw.(string)

		logger.Info("incoming request",
			slog.String("req_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("query", c.Context().QueryArgs().String()),
		)

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		if err != nil {
			logger.Error("request finished with error",
				slog.String("req_id", reqID),
				slog.Int("status", status),
				slog.String("duration", latency.String()),
				slog.String("error", err.Error()),
			)

			return err
		}

		logger.Info("request finished",
			slog.String("req_id", reqID),
			slog.Int("status", status),
			slog.String("duration", latency.String()),
		)

		return nil
	}
}
