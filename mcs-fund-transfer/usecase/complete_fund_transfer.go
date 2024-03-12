package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
)

type completeFundTransferHandler struct {
	logger    logger.Logger
	fData     data.FundTransferData
	fWorkflow outbound.FundTransferWorkflow
}

func NewCompleteFundTransferHandler(log logger.Logger, fData data.FundTransferData, fWorkflow outbound.FundTransferWorkflow) domain.CompleteFundTransferHandler {
	return &completeFundTransferHandler{
		logger:    log,
		fData:     fData,
		fWorkflow: fWorkflow,
	}
}

// Handle implements domain.CompleteFundTransferHandler.
func (h *completeFundTransferHandler) Handle(ctx context.Context, req *domain.CompleteFundTransferRequest) (*domain.CompleteFundTransferResponse, error) {
	trans, err := h.fData.GetFundTransferTransactionByTransNo(ctx, req.TransNo)
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

	// signal that the transaction completed successfully
	err = h.fWorkflow.SignalFundTransferVerifiedOTP(ctx, trans)
	if err != nil {
		return nil, err
	}

	h.logger.Infof(ctx, "Signaled transaction completed successfully. Workflow ID: %v", trans.WorflowId)

	return &domain.CompleteFundTransferResponse{
		CrefNum: trans.CRefNum,
		TransNo: req.TransNo,
	}, nil
}
