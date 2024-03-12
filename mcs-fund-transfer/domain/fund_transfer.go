package domain

import (
	"context"
	"time"
)

// Start fund transfer
type StartFundTransferRequest struct {
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	Amount      int64  `json:"amount"`
}

type StartFundTransferResponse struct {
	CRefNum    string `json:"cRefNum"`
	WorkflowId string `json:"workflowId"`
}

type StartFundTransferHandler interface {
	Handle(context.Context, *StartFundTransferRequest) (*StartFundTransferResponse, error)
}

// Generate OTP
type GenerateFundTransferOTPRequest struct {
	CrefNum string `json:"cRefNum"`
}

type GenerateFundTransferOTPResponse struct {
}

type GenerateFundTransferOTPHandler interface {
	Handle(context.Context, *GenerateFundTransferOTPRequest) (*GenerateFundTransferOTPResponse, error)
}

// Verify OTP
type VerifyFundTransferOTPRequest struct {
	CrefNum string `json:"cRefNum"`
	OTP     string `json:"otp"`
}

type VerifyFundTransferOTPResponse struct {
	Success string `json:"success"`
}

type VerifyFundTransferOTPHandler interface {
	Handle(ctx context.Context, request *VerifyFundTransferOTPRequest) (*VerifyFundTransferOTPResponse, error)
}

// Execute Fund transfer
type ExecuteFundTransferRequest struct {
	CRefNum string `json:"cRefNum"`
}

type ExecuteFundTransferResponse struct {
}

type ExecuteFundTransferHandler interface {
	Handle(context.Context, *ExecuteFundTransferRequest) (*ExecuteFundTransferResponse, error)
}

type CompleteFundTransferRequest struct {
	CRefNum    string    `json:"cRefNum"`
	TransferAt time.Time `json:"transferAt"`
}

type CompleteFundTransferResponse struct {
}

type CompleteFundTransferHandler interface {
	Handle(context.Context, *CompleteFundTransferRequest) (*CompleteFundTransferResponse, error)
}
