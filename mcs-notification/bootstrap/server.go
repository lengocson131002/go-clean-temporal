package bootstrap

import "github.com/lengocson131002/go-clean-core/config"

type ServerConfig struct {
	Name               string
	AppVersion         string
	HttpPort           int
	GrpcPort           int
	BaseURI            string
	GrRunningThreshold int // threshold for goroutines are running (which could indicate a resource leak).
	GcPauseThresholdMs int // threshold threshold garbage collection pause exceeds. (Millisecond)
	EnvFilePath        string
}

func GetServerConfig(cfg config.Configure) *ServerConfig {
	name := cfg.GetString("APP_NAME")
	version := cfg.GetString("APP_VERSION")
	httpPort := cfg.GetInt("APP_HTTP_PORT")
	grpcPort := cfg.GetInt("APP_GRPC_PORT")
	baseUrl := cfg.GetString("APP_BASE_URL")
	grRunningThreshold := cfg.GetInt("APP_GR_RUNNING_THRESHOLD")
	gcMaxPauseThresholdms := cfg.GetInt("APP_GC_PAUSE_THRESHOLD_MS")

	return &ServerConfig{
		Name:               name,
		AppVersion:         version,
		HttpPort:           httpPort,
		GrpcPort:           grpcPort,
		BaseURI:            baseUrl,
		GrRunningThreshold: grRunningThreshold,
		GcPauseThresholdMs: gcMaxPauseThresholdms,
		EnvFilePath:        "./.env",
	}
}
