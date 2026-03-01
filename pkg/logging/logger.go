package logging

import (
	"minigate/pkg/config"

	"go.uber.org/zap"
)

func New(serviceName string) *zap.Logger {
	level := config.GetEnv("LOG_LEVEL", "info")
	cfg := zap.NewProductionConfig()
	if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewExample()
	}
	return logger.With(zap.String("service", serviceName))
}
