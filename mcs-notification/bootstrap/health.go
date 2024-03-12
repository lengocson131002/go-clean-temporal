package bootstrap

import (
	"time"

	health "github.com/lengocson131002/go-clean-core/health"
)

func NewHealthChecker(cfg *ServerConfig) health.HealthChecker {
	// Init health
	healthChecker := health.NewHealthChecker(cfg.Name, cfg.AppVersion)

	// check Garbage Collector
	gcChecker := health.NewGarbageCollectionMaxChecker(time.Millisecond * time.Duration(cfg.GcPauseThresholdMs))
	healthChecker.AddLivenessCheck("garbage collector check", gcChecker)

	// check Goroutine
	grChecker := health.NewGoroutineChecker(cfg.GrRunningThreshold)
	healthChecker.AddLivenessCheck("goroutine checker", grChecker)

	// check env file
	envFileChecker := health.NewEnvChecker(cfg.EnvFilePath)
	healthChecker.AddReadinessCheck("env file checker", envFileChecker)

	// check network
	pingChecker := health.NewPingChecker("http://google.com", "GET", time.Millisecond*time.Duration(200), nil, nil)
	healthChecker.AddReadinessCheck("ping check", pingChecker)

	return healthChecker
}
