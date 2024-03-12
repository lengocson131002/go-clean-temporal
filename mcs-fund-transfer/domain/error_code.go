package domain

import (
	"fmt"
	"net/http"

	"github.com/lengocson131002/go-clean-core/errors"
)

var (
	// DOMAIN CUSTOM ERROR
	ErrorTransactionNotFound = &errors.DomainError{
		Status:  http.StatusBadRequest,
		Code:    "100",
		Message: "Transaction not found",
	}

	ErrorInvalidOTP = &errors.DomainError{
		Status:  http.StatusBadRequest,
		Code:    "101",
		Message: "Invalid OTP",
	}
)

func ErrorFailedToExecuteTransaction(detail string) *errors.DomainError {
	return &errors.DomainError{
		Status:  http.StatusBadRequest,
		Code:    "103",
		Message: fmt.Sprintf("Failed to execute transaction payment. Details: %s", detail),
	}
}
