package data

import (
	"context"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
)

type FundTransferData interface {
	SaveFundTransferTransaction(context.Context, *domain.FunTransferTransaction) error
	SaveFundTransferOTP(context.Context, *domain.FundTransferOTP) error
	GetFundTransferTransaction(ctx context.Context, cRefNum string) (*domain.FunTransferTransaction, error)
	GetFundTransferTransactionByTransNo(ctx context.Context, transNo string) (*domain.FunTransferTransaction, error)
	GetFundTransferOTP(ctx context.Context, cRefNum string, otp string) (*domain.FundTransferOTP, error)
}
