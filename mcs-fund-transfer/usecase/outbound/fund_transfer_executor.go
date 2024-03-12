package outbound

import (
	"context"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
)

type ExecutePaymentResponse struct {
	Success bool
	Detail  string
}

type FundTransferExecutor interface {
	ExecutePayment(ctx context.Context, req *domain.FunTransferTransaction) (*ExecutePaymentResponse, error)
}
