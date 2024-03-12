package domain

import "time"

type FundTransferStatus int

const (
	TransactionStarted FundTransferStatus = iota
	TransactionVerified
	TransactionProcessing
	TransactionSucceeded
)

type FunTransferTransaction struct {
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

type FundTransferOTP struct {
	OTP       string    `json:"otp"`
	CRefNum   string    `json:"cRefNum"`
	CreatedAt time.Time `json:"createdAt"`
	Verified  bool      `json:"verified"`
}
