package bad

import (
	"log/slog"

	"go.uber.org/zap"
)

const cMsg = "Starting server"

func bad() {
	slog.Info("Starting server") // want "log message must start with a lowercase letter"

	slog.Info("запуск") // want "log message must be in English" "log message must not contain special characters or emojis"

	slog.Info("server started!🚀") // want "log message must be in English" "log message must not contain special characters or emojis"

	slog.Info("auth_token") // want "log message must not contain potentially sensitive data"

	slog.Info(cMsg) // want "log message must start with a lowercase letter"

	slog.Info("Auth" + " token") // want "log message must start with a lowercase letter" "log message must not contain potentially sensitive data"

	logger, _ := zap.NewProduction()

	logger.Info("Starting server") // want "log message must start with a lowercase letter"

}
