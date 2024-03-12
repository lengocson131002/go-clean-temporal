package data

import "context"

type AccountBalanceResponse struct {
	Currency        string
	OpenActualBal   int64
	OnlineActualBal int64
	WorkingBalance  int64
}

type AccountData interface {
	GetBalance(cxt context.Context, accNumber string) (*AccountBalanceResponse, error)
}
