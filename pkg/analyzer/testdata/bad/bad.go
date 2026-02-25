package bad

import (
	"log/slog"

	"go.uber.org/zap"
)

func bad() {
	slog.Info("Starting server on port 8080") // want "log message must start with a lowercase letter"
	slog.Error("запуск сервера")              // want "log message must be in English" "log message must not contain special characters or emojis"
	slog.Info("server started!🚀")             // want "log message must be in English" "log message must not contain special characters or emojis"
	slog.Error("api_key=123")                 // want "log message must not contain special characters or emojis" "log message must not contain potentially sensitive data"

	logger, _ := zap.NewProduction()
	logger.Info("Starting server")    // want "log message must start with a lowercase letter"
	logger.Info("ошибка подключения") // want "log message must be in English" "log message must not contain special characters or emojis"
	logger.Info("server started!!!")  // want "log message must not contain special characters or emojis"
	logger.Info("api_key=123")        // want "log message must not contain special characters or emojis" "log message must not contain potentially sensitive data"
}
