package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean-core/util"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
)

type queryFundTransferUserCase struct {
	fData data.FundTransferData
}

func NewQueryFundTransferData(fData data.FundTransferData) domain.QueryFundTransferHandler {
	return &queryFundTransferUserCase{
		fData: fData,
	}
}

func (h *queryFundTransferUserCase) Handle(ctx context.Context, request *domain.QueryFundTransferRequest) (*domain.QueryFundTransferResponse, error) {
	trans, err := h.fData.GetFundTransferTransaction(ctx, request.CrefNum)
	if err != nil {
		return nil, err
	}

	if trans == nil {
		return nil, domain.ErrorTransactionNotFound
	}

	var res domain.QueryFundTransferResponse
	err = util.MapStruct(trans, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
