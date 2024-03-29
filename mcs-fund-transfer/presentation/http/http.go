package http

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	healthchecks "github.com/lengocson131002/go-clean-core/health"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/bootstrap"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/http/controller"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/http/handler"
)

type HttpServer struct {
	app                    *fiber.App
	cfg                    *bootstrap.ServerConfig
	logger                 logger.Logger
	healhChecker           healthchecks.HealthChecker
	fundTransferController *controller.FundTransferController
}

// @title  GOLANG TEMPORAL DEMO
// @version 1.0
// @description GOLANG TEMPORAL DEMO
// @termsOfService http://swagger.io/terms/
// @contact.name LNS
// @contact.email leson131002@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func NewHttpServer(
	cfg *bootstrap.ServerConfig,
	logger logger.Logger,
	healhChecker healthchecks.HealthChecker,
) *HttpServer {
	return &HttpServer{
		cfg:          cfg,
		logger:       logger,
		healhChecker: healhChecker,
	}
}

func (s *HttpServer) Start(ctx context.Context, opts ...ServerOption) error {
	s.app = fiber.New(fiber.Config{
		ErrorHandler: handler.CustomErrorHandler,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	// configs options
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			if err := s.app.Shutdown(); err != nil {
				s.logger.Errorf(ctx, "Failed to shutdown http server: %v", err)
			}
			s.logger.Info(ctx, "Stop HTTP Server")
		}(ctx)

		<-ctx.Done()
	}()

	hPort := s.cfg.HttpPort
	s.logger.Infof(ctx, "Start HTTP server at port: %v", hPort)
	if err := s.app.Listen(fmt.Sprintf(":%v", hPort)); err != nil {
		s.logger.Errorf(ctx, "Failed to start http server: %v ", err)
		return err
	}

	return nil
}

// OPTIONS
type ServerOption func(*HttpServer) error

func WithLoggings() ServerOption {
	return func(s *HttpServer) error {
		s.app.Use(fiberLog.New(fiberLog.Config{
			Next:         nil,
			Done:         nil,
			Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stdout,
		}))
		return nil
	}
}

func WithSwagger() ServerOption {
	return func(s *HttpServer) error {
		s.app.Get("/swagger/*", swagger.HandlerDefault)
		return nil
	}
}

func WithHealthCheck() ServerOption {
	return func(s *HttpServer) error {
		s.app.Get("/liveliness", func(c *fiber.Ctx) error {
			result := s.healhChecker.LivenessCheck()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})

		s.app.Get("/readiness", func(c *fiber.Ctx) error {
			result := s.healhChecker.RedinessCheck()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})
		return nil
	}
}

func WithTracing() ServerOption {
	return func(s *HttpServer) error {
		s.app.Use(otelfiber.Middleware())
		return nil
	}
}

func WithMetrics() ServerOption {
	return func(s *HttpServer) error {
		prometheus := fiberprometheus.New(s.cfg.Name)
		prometheus.RegisterAt(s.app, "/metrics")
		s.app.Use(prometheus.Middleware)
		return nil
	}
}

func WithRoutes() ServerOption {
	return func(s *HttpServer) error {
		var (
			fApp  = s.app
			fConn = s.fundTransferController
		)

		appV1 := fApp.Group("/api/v1")

		appV1.Post("/fund-transfer/start", fConn.FundTransfer)
		appV1.Post("/fund-transfer/verify-otp", fConn.VerifyOTP)
		appV1.Get("/fund-transfer/query", fConn.QueryFundTransfer)

		s.logger.Infof(context.Background(), "[Fiber] Register v1 routes")
		return nil
	}
}
