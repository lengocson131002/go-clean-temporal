package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domainErrors "github.com/lengocson131002/go-clean-core/errors"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/metrics/prome"
	"github.com/lengocson131002/go-clean-core/pipeline"
	"github.com/lengocson131002/go-clean-core/trace"
	"github.com/lengocson131002/go-clean-core/util"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterPipelineHandlers(
	startFundTransferHandler domain.StartFundTransferHandler,
	generateFundTransferOTPHandler domain.GenerateFundTransferOTPHandler,
	verifyOTPHandler domain.VerifyFundTransferOTPHandler,
	executeFundTransferHandler domain.ExecuteFundTransferHandler,
) {
	pipeline.RegisterRequestHandler(startFundTransferHandler)
	pipeline.RegisterRequestHandler(generateFundTransferOTPHandler)
	pipeline.RegisterRequestHandler(verifyOTPHandler)
	pipeline.RegisterRequestHandler(executeFundTransferHandler)
}

// ERROR HANDLING FOR RECOVERING FROM PANIC
type ErrorHandlingBehavior struct {
	logger logger.Logger
}

func NewErrorHandlingBehavior(logger logger.Logger) *ErrorHandlingBehavior {
	return &ErrorHandlingBehavior{
		logger: logger,
	}
}

func (b *ErrorHandlingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (res interface{}, err error) {
	// TODO: recover from error panic to prevent stop application
	defer func() {
		if r := recover(); r != nil {
			b.logger.Errorf(ctx, "Recovered from panic: %v", r)
			err = fmt.Errorf("internal server error")
		}
	}()
	response, err := next(ctx)
	return response, err
}

// TRACING
type RequestTracingBehavior struct {
	logger logger.Logger
	tracer trace.Tracer
}

func NewTracingBehavior(logger logger.Logger, tracer trace.Tracer) *RequestTracingBehavior {
	return &RequestTracingBehavior{
		logger: logger,
		tracer: tracer,
	}
}

func (b *RequestTracingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (interface{}, error) {
	reqType := util.GetType(request)
	opName := fmt.Sprintf("request pipeline - %s", reqType)

	// tracing request
	var res interface{}
	ctx, finish := b.tracer.StartInternalTrace(ctx, opName, trace.WithInternalRequest(request))
	defer finish(ctx, trace.WithInternalResponse(res))

	res, err := next(ctx)

	return res, err
}

// METRICS
type RequestMetricBehavior struct {
	logger   logger.Logger
	metricer *prome.PrometheusMetricer
	cfg      *ServerConfig
}

func NewMetricBehavior(logger logger.Logger, metricer *prome.PrometheusMetricer, cfg *ServerConfig) *RequestMetricBehavior {
	return &RequestMetricBehavior{
		logger:   logger,
		metricer: metricer,
		cfg:      cfg,
	}
}

func (b *RequestMetricBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	reqType := util.GetType(request)

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		b.metricer.RequestSummary.WithLabelValues(b.cfg.Name, reqType).Observe(us)
		b.metricer.RequestHistogram.WithLabelValues(b.cfg.Name, reqType).Observe(v)
	}))

	defer func() {
		var businessErr *domainErrors.DomainError
		// mark a request success as if there is no error happened or the error is business error
		if err == nil || errors.As(err, &businessErr) {
			b.metricer.RequestTotalCounter.WithLabelValues(b.cfg.Name, reqType, "success").Inc()
		} else {
			b.metricer.RequestTotalCounter.WithLabelValues(b.cfg.Name, reqType, "failure").Inc()
		}

		// stop timer
		timer.ObserveDuration()
	}()

	response, err = next(ctx)
	return response, err
}

// LOGGING
type RequestLoggingBehavior struct {
	logger logger.Logger
}

func NewRequestLoggingBehavior(logger logger.Logger) *RequestLoggingBehavior {
	return &RequestLoggingBehavior{
		logger: logger,
	}
}

func (b *RequestLoggingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	isSuccess := true
	start := time.Now()

	defer func() {
		if err != nil {
			isSuccess = false
		}

		var (
			requestJson, _  = json.Marshal(request)
			responseJson, _ = json.Marshal(response)
			errJson         string
		)

		if err != nil {
			errJson = err.Error()
		}

		b.logger.Infof(ctx, "[Request Pipeline] Success: %t - Request: %s - Response: %s - Error: %s - Duration: %dms",
			isSuccess,
			string(requestJson),
			string(responseJson),
			errJson,
			time.Since(start).Milliseconds())
	}()

	response, err = next(ctx)

	return response, err
}

func RegisterPipelineBehaviors(
	requestLoggingBehavior *RequestLoggingBehavior,
	requestTracingBehavior *RequestTracingBehavior,
	requestMetricBehavior *RequestMetricBehavior,
	errorHandlingBehavior *ErrorHandlingBehavior,

) {
	pipeline.RegisterRequestPipelineBehaviors(
		requestTracingBehavior,
		requestLoggingBehavior,
		requestMetricBehavior,
		errorHandlingBehavior,
	)

}
