package gprc

import (
	"context"
	"fmt"
	"net"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/bootstrap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	cfg    *bootstrap.ServerConfig
	logger logger.Logger
	srv    *grpc.Server
}

type GrpcServerOption func(*GrpcServer) error

func NewGrpcServer(cfg *bootstrap.ServerConfig, logger logger.Logger) *GrpcServer {
	return &GrpcServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *GrpcServer) Start(ctx context.Context, opts ...GrpcServerOption) error {
	network, gPort := "tcp", s.cfg.GrpcPort
	lis, err := net.Listen(network, fmt.Sprintf("localhost:%d", gPort))

	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func() {
			s.srv.GracefulStop()
			s.logger.Info(ctx, "Stop GRPC Server")
		}()
		<-ctx.Done()
	}()

	s.logger.Infof(ctx, "Start GRPC server at port: %v", gPort)
	if err := s.srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve GRPC %w", err)
	}

	return nil
}

func (g *GrpcServer) WithT24Server() GrpcServerOption {
	return func(s *GrpcServer) error {
		// tSrv := server.NewT24AccountServer(s.logger)
		// pb.RegisterT24AccountServiceServer(gSrv, tSrv)
		return nil
	}
}
