package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
)

type verifyOTPHandler struct {
	fData     data.FundTransferData
	fWorkflow outbound.FundTransferWorkflow
	logger    logger.Logger
}

func NewVerifyOTPHandler(
	fData data.FundTransferData,
	fWorkflow outbound.FundTransferWorkflow,
	logger logger.Logger,
) domain.VerifyFundTransferOTPHandler {
	return &verifyOTPHandler{
		fData:     fData,
		fWorkflow: fWorkflow,
		logger:    logger,
	}
}

func (h *verifyOTPHandler) Handle(ctx context.Context, request *domain.VerifyFundTransferOTPRequest) (*domain.VerifyFundTransferOTPResponse, error) {
	otp, err := h.fData.GetFundTransferOTP(ctx, request.CrefNum, request.OTP)
	if err != nil {
		return nil, err
	}

	if otp == nil || otp.Verified {
		return nil, domain.ErrorInvalidOTP
	}

	otp.Verified = true
	err = h.fData.SaveFundTransferOTP(ctx, otp)
	if err != nil {
		return nil, err
	}

	trans, err := h.fData.GetFundTransferTransaction(ctx, request.CrefNum)
	if err != nil {
		return nil, err
	}

	if trans == nil {
		return nil, domain.ErrorTransactionNotFound
	}

	trans.Status = domain.TransactionVerified
	err = h.fData.SaveFundTransferTransaction(ctx, trans)
	if err != nil {
		h.logger.Errorf(ctx, "failed to update transaction status: %v", err)
		return nil, err
	}

	// TODO: signal workflow
	err = h.fWorkflow.SignalFundTransferVerifiedOTP(ctx, trans)
	if err != nil {
		return nil, err
	}

	h.logger.Infof(ctx, "Signaled fund tranfer verified OTP. Workflow ID: %s", trans.WorflowId)

	return &domain.VerifyFundTransferOTPResponse{
		Success: true,
	}, nil
}
