package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
)

type completeFundTransferHandler struct {
	logger logger.Logger
	fData  data.FundTransferData
}

func NewCompleteFundTransferHandler(log logger.Logger, fData data.FundTransferData) domain.CompleteFundTransferHandler {
	return &completeFundTransferHandler{
		logger: log,
		fData:  fData,
	}
}

// Handle implements domain.CompleteFundTransferHandler.
func (h *completeFundTransferHandler) Handle(ctx context.Context, req *domain.CompleteFundTransferRequest) (*domain.CompleteFundTransferResponse, error) {
	trans, err := h.fData.GetFundTransferTransaction(ctx, req.CRefNum)
	if err != nil {
		return nil, err
	}

	if trans == nil || trans.Status != domain.TransactionProcessing {
		return nil, domain.ErrorTransactionNotFound
	}

	trans.Status = domain.TransactionSucceeded
	trans.TransferAt = &req.TransferAt

	err = h.fData.SaveFundTransferTransaction(ctx, trans)
	if err != nil {
		return nil, err
	}

	return &domain.CompleteFundTransferResponse{}, nil
}
