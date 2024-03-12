package broker

import (
	"context"
	"sync"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/transport/broker"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
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
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				_, err := b.broker.Subscribe(TopicRequestGenerateOTP, func(e broker.Event) error {
					return HandleBrokerEvent[*domain.GenerateFundTransferOTPRequest, *domain.GenerateFundTransferOTPResponse](b.broker, e, WithReplyTopic(TopicReplyGenerateOTP))
				})
				if err != nil {
					cancel()
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				_, err := b.broker.Subscribe(TopicRequestExecuteFundTransfer, func(e broker.Event) error {
					return HandleBrokerEvent[*domain.ExecuteFundTransferRequest, *domain.ExecuteFundTransferResponse](b.broker, e, WithReplyTopic(TopicReplyExecuteFundTransfer))
				})
				if err != nil {
					cancel()
				}
			}

		}()

		wg.Wait()
		return ctx.Err()
	}
}
