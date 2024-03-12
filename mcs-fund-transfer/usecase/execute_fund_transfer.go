package usecase

import (
	"context"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
)

type executeFundTransferHandler struct {
	fData     data.FundTransferData
	fExecutor outbound.FundTransferExecutor
}

func NewExecuteFundTransferHandler(
	fData data.FundTransferData,
	fExecutor outbound.FundTransferExecutor,

) domain.ExecuteFundTransferHandler {
	return &executeFundTransferHandler{
		fData:     fData,
		fExecutor: fExecutor,
	}
}

func (h *executeFundTransferHandler) Handle(ctx context.Context, req *domain.ExecuteFundTransferRequest) (*domain.ExecuteFundTransferResponse, error) {
	trans, err := h.fData.GetFundTransferTransaction(ctx, req.CRefNum)
	if err != nil {
		return nil, err
	}

	if trans == nil || (trans.Status != domain.TransactionVerified) {
		return nil, domain.ErrorTransactionNotFound
	}

	res, err := h.fExecutor.ExecutePayment(ctx, trans)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, domain.ErrorFailedToExecuteTransaction(res.Detail)
	}

	trans.Status = domain.TransactionProcessing
	trans.TransNo = res.TransNo
	err = h.fData.SaveFundTransferTransaction(ctx, trans)
	if err != nil {
		return nil, err
	}

	return &domain.ExecuteFundTransferResponse{}, nil
}
