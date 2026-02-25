package good

import (
	"log/slog"

	"go.uber.org/zap"
)

func good() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Info("server started")
	slog.Info("user authenticated successfully")

	logger, _ := zap.NewProduction()
	logger.Info("starting server on port 8080")
	logger.Error("failed to connect to database")
	logger.Warn("something went wrong")
	logger.Info("api request completed")
}
