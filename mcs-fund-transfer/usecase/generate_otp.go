package usecase

import (
	"context"
	"time"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/pkg/utils"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
)

type generateFundTransferOTPHandler struct {
	fData  data.FundTransferData
	logger logger.Logger
}

func NewGenerateFundTransferOTPHandler(
	fData data.FundTransferData,
	logger logger.Logger,
) domain.GenerateFundTransferOTPHandler {
	return &generateFundTransferOTPHandler{
		fData:  fData,
		logger: logger,
	}
}

func (h *generateFundTransferOTPHandler) Handle(ctx context.Context, request *domain.GenerateFundTransferOTPRequest) (*domain.GenerateFundTransferOTPResponse, error) {
	// check transaction
	trans, err := h.fData.GetFundTransferTransaction(ctx, request.CrefNum)
	if err != nil {
		return nil, err
	}

	opt, err := utils.GenerateOTP(6)
	if err != nil {
		return nil, err
	}

	err = h.fData.SaveFundTransferOTP(ctx, &domain.FundTransferOTP{
		OTP:       opt,
		CRefNum:   trans.CRefNum,
		CreatedAt: time.Now(),
		Verified:  false,
	})

	if err != nil {
		return nil, err
	}

	h.logger.Infof(ctx, "generated fund transfer OTP. workflow ID: %v", trans.WorflowId)

	return nil, nil
}
