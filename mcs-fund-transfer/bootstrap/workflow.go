package bootstrap

import (
	"context"

	"github.com/lengocson131002/go-clean-core/config"
	"github.com/lengocson131002/go-clean-core/logger"
	"go.temporal.io/sdk/client"
)

func NewWorkflowClient(cfg config.Configure, logger logger.Logger) client.Client {
	var (
		hostPort  = cfg.GetString("TEMPORAL_HOST_PORT")
		namespace = cfg.GetString("TEMPORAL_NAMESPACE")
	)

	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	})

	if err != nil {
		panic(err)
	}

	logger.Info(context.TODO(), "Connected to Elasticsearch")
	return c
}
