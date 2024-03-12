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
	Success bool `json:"success"`
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
	TransNo    string    `json:"transNo"`
	TransferAt time.Time `json:"transferAt"`
}

type CompleteFundTransferResponse struct {
	CrefNum string `json:"crefNum"`
	TransNo string `json:"transNo"`
}

type CompleteFundTransferHandler interface {
	Handle(context.Context, *CompleteFundTransferRequest) (*CompleteFundTransferResponse, error)
}

type QueryFundTransferRequest struct {
	CrefNum string `json:"cRefNum"`
}

type QueryFundTransferResponse struct {
	WorflowId   string             `json:"worflowId"`
	FromAccount string             `json:"fromAccount"`
	ToAccount   string             `json:"toAccount"`
	Amount      int64              `json:"amount"`
	CRefNum     string             `json:"cRefNum"`
	CreatedAt   time.Time          `json:"createdAt"`
	TransferAt  *time.Time         `json:"transferAt"`
	Status      FundTransferStatus `json:"status"`
	TransNo     string             `json:"transNo"`
}

type QueryFundTransferHandler interface {
	Handle(context.Context, *QueryFundTransferRequest) (*QueryFundTransferResponse, error)
}
