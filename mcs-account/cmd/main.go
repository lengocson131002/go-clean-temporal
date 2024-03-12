package main

import (
	"context"
	"time"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-account/bootstrap"
	"github.com/lengocson131002/mcs-account/infras/data"
	"github.com/lengocson131002/mcs-account/presentation/broker"
	"github.com/lengocson131002/mcs-account/presentation/http"
	"github.com/lengocson131002/mcs-account/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("main",
	fx.Provide(bootstrap.GetLogger),
	fx.Provide(bootstrap.GetConfigure),
	fx.Provide(bootstrap.GetServerConfig),
	fx.Provide(bootstrap.GetValidator),
	fx.Provide(bootstrap.GetTracer),
	fx.Provide(bootstrap.GetKafkaBroker),
	fx.Provide(bootstrap.NewHealthChecker),
	fx.Provide(bootstrap.NewElasticSearchClient),

	// PIPELINE
	fx.Provide(bootstrap.NewMetricBehavior),
	fx.Provide(bootstrap.NewTracingBehavior),
	fx.Provide(bootstrap.NewRequestLoggingBehavior),
	fx.Provide(bootstrap.NewErrorHandlingBehavior),
	fx.Provide(usecase.NewCheckBalanceHandler),

	// INFRAS
	fx.Provide(data.NewAccountData),

	// SERVERS
	fx.Provide(http.NewHttpServer),
	fx.Provide(broker.NewBrokerServer),

	fx.Provide(bootstrap.GetPrometheusMetricer),
	fx.Invoke(bootstrap.RegisterPipelineBehaviors),
	fx.Invoke(bootstrap.RegisterPipelineHandlers),
)

func main() {
	// Dependencies injection using FX package
	fx.New(
		Module,
		fx.Invoke(run),
	).Run()

}

func run(lc fx.Lifecycle, httpServer *http.HttpServer, brokerServer *broker.BrokerServer, log logger.Logger, conf *bootstrap.ServerConfig, shutdowner fx.Shutdowner) {
	gCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			errChan := make(chan error)

			// start HTTP server
			go func() {
				err := httpServer.Start(gCtx,
					http.WithHealthCheck(),
					http.WithLoggings(),
					http.WithTracing(),
					http.WithMetrics())

				if err != nil {
					log.Fatal(ctx, "Failed to start HTTP server: %s", err)
					errChan <- err
					cancel()
					shutdowner.Shutdown()
				}
			}()

			// start BROKER server
			go func() {
				if err := brokerServer.Start(gCtx, broker.SubscribeHandler()); err != nil {
					log.Fatalf(ctx, "Failed to start Broker server: %s", err)
					errChan <- err
					cancel()
					shutdowner.Shutdown()
				}
			}()

			select {
			case err := <-errChan:
				return err
			case <-time.After(100 * time.Millisecond):
				return nil
			}

		},
		OnStop: func(ctx context.Context) error {
			cancel()
			select {
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		},
	})
}
