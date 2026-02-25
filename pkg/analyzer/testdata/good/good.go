package good

import (
	"log/slog"

	"go.uber.org/zap"
)

func good() {
	slog.Info("starting server")
	slog.Error("failed to connect to database")

	logger, _ := zap.NewProduction()
	logger.Info("starting server")
	logger.Error("failed to connect to database")
}
