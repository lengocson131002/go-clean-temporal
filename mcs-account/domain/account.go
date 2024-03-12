package domain

import "context"

type CheckBalanceRequest struct {
	Account string `json:"account"`
}

type CheckBalanceResponse struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type CheckBalanceHandler interface {
	Handle(ctx context.Context, request *CheckBalanceRequest) (*CheckBalanceResponse, error)
}
