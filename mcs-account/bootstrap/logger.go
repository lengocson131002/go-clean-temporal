package bootstrap

import (
	"github.com/lengocson131002/go-clean-core/config"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/logger/logrus"
	"github.com/lengocson131002/go-clean-core/trace"
)

func GetLogger(c config.Configure, tracer trace.Tracer) logger.Logger {
	levelStr := c.GetString("LOG_LEVEL")
	level, err := logger.GetLevel(levelStr)
	if err != nil {
		level = logger.InfoLevel
	}
	return logrus.NewLogrusLogger(
		logger.WithLevel(level),
		logger.WithTracer(tracer),
	)
}
