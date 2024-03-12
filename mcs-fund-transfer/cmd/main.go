package main

import (
	"context"
	"time"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/bootstrap"
	"github.com/lengocson131002/mcs-fund-transfer/infras/data"
	"github.com/lengocson131002/mcs-fund-transfer/infras/outbound"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/broker"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/cron"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/http"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/http/controller"
	"github.com/lengocson131002/mcs-fund-transfer/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("main",
	// Bootstrap
	fx.Provide(bootstrap.GetLogger),
	fx.Provide(bootstrap.GetConfigure),
	fx.Provide(bootstrap.GetServerConfig),
	fx.Provide(bootstrap.GetValidator),
	fx.Provide(bootstrap.GetTracer),
	fx.Provide(bootstrap.GetKafkaBroker),
	fx.Provide(bootstrap.NewHealthChecker),
	fx.Provide(bootstrap.NewWorkflowClient),
	fx.Provide(bootstrap.NewElasticSearchClient),
	fx.Provide(bootstrap.GetPrometheusMetricer),

	// Presentaion
	fx.Provide(http.NewHttpServer),
	fx.Provide(broker.NewBrokerServer),
	fx.Provide(controller.NewFundTransferController),
	fx.Provide(cron.NewCronServer),

	// Usecase
	fx.Provide(bootstrap.NewMetricBehavior),
	fx.Provide(bootstrap.NewTracingBehavior),
	fx.Provide(bootstrap.NewRequestLoggingBehavior),
	fx.Provide(bootstrap.NewErrorHandlingBehavior),
	fx.Provide(usecase.NewStartFundTransferHandler),
	fx.Provide(usecase.NewVerifyOTPHandler),
	fx.Provide(usecase.NewExecuteFundTransferHandler),
	fx.Provide(usecase.NewGenerateFundTransferOTPHandler),
	fx.Invoke(bootstrap.RegisterPipelineBehaviors),
	fx.Invoke(bootstrap.RegisterPipelineHandlers),

	// INFRAS
	fx.Provide(outbound.NewFundTransferWorkflow),
	fx.Provide(data.NewFundTransferData),
	fx.Provide(outbound.NewFundTransferExecutor),
)

func main() {
	// Dependencies injection using FX package
	fx.New(
		Module,
		fx.Invoke(run),
	).Run()

}

func run(
	lc fx.Lifecycle,
	httpServer *http.HttpServer,
	brokerServer *broker.BrokerServer,
	cronServer *cron.CronServer,
	log logger.Logger,
	conf *bootstrap.ServerConfig,
	shutdowner fx.Shutdowner) {
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
					http.WithMetrics(),
					http.WithV1Routes(),
				)

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

			// start CRON server
			go func() {
				if err := cronServer.Start(gCtx, cron.WithCompleteFundTransferBackground()); err != nil {
					log.Fatalf(ctx, "Failed to start cron server: %s", err)
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
