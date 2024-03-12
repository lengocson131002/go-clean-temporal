package broker

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

type BrokerServer struct {
	broker broker.Broker
	logger logger.Logger
}

func NewBrokerServer(broker broker.Broker, logger logger.Logger) *BrokerServer {
	return &BrokerServer{
		broker: broker,
		logger: logger,
	}
}

type BrokerServerOption func(*BrokerServer) error

func (s *BrokerServer) Start(ctx context.Context, opts ...BrokerServerOption) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			if err := s.broker.Disconnect(); err != nil {
				s.logger.Errorf(ctx, "Failed to shutdown broker server: %v", err)
			}
			s.logger.Info(ctx, "Stop Broker Server")
		}(ctx)

		<-ctx.Done()
	}()

	return nil
}

func SubscribeHandler() BrokerServerOption {
	return func(b *BrokerServer) error {
		_, err := b.broker.Subscribe("", func(e broker.Event) error {
			return nil
		})
		return err
	}
}
