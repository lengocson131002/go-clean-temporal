package cron

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/pipeline"
	"github.com/lengocson131002/go-clean-core/transport/broker"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
)

const (
	TopicAccountStatement = "OCBDWDAILY.TODAY_ACCOUNT_STATEMENT"
)

type CronServer struct {
	logger logger.Logger
	broker broker.Broker
}

type CronServerOption func(*CronServer) error

func NewCronServer(logger logger.Logger, broker broker.Broker) *CronServer {
	return &CronServer{
		logger: logger,
		broker: broker,
	}
}

func (s *CronServer) Start(ctx context.Context, opts ...CronServerOption) error {
	// configs options
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			s.logger.Info(ctx, "Stop Cron Server")
		}(ctx)
		<-ctx.Done()
	}()

	s.logger.Infof(ctx, "Started Cron Server")
	return nil
}

func WithCompleteFundTransferBackground() CronServerOption {
	return func(cs *CronServer) error {
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
				_, err := cs.broker.Subscribe(TopicAccountStatement, func(e broker.Event) error {
					if e.Message() == nil || len(e.Message().Body) == 0 {
						// ignore
						return nil
					}

					body := e.Message().Body
					var request KafkaAccountStatementWrapper
					err := json.Unmarshal(body, &request)
					if err != nil {
						return broker.InvalidDataFormatError{}
					}

					if request.OpType == "I" && request.After != nil {
						pReq := domain.CompleteFundTransferRequest{
							TransNo:    request.After.TransNo,
							TransferAt: time.Now(),
						}
						ctx := context.Background()
						res, err := pipeline.Send[*domain.CompleteFundTransferRequest, *domain.CompleteFundTransferResponse](ctx, &pReq)
						if err != nil {
							cs.logger.Errorf(ctx, "Complete fund transfer failed: %v", err.Error())
						} else {
							cs.logger.Errorf(ctx, "Completed fund transfer: %v", res)
						}
					}
					return nil
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
