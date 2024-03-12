package usecase

import (
	"context"

	"github.com/lengocson131002/mcs-account/domain"
	"github.com/lengocson131002/mcs-account/usecase/data"
)

type checkBalanceHandler struct {
	accDb data.AccountData
}

func NewCheckBalanceHandler(
	accDb data.AccountData,
) domain.CheckBalanceHandler {
	return &checkBalanceHandler{
		accDb: accDb,
	}
}

func (h *checkBalanceHandler) Handle(ctx context.Context, request *domain.CheckBalanceRequest) (*domain.CheckBalanceResponse, error) {
	balRes, err := h.accDb.GetBalance(ctx, request.Account)
	if err != nil {
		return nil, err
	}

	return &domain.CheckBalanceResponse{
		Balance:  balRes.WorkingBalance,
		Currency: balRes.Currency,
	}, nil

}
