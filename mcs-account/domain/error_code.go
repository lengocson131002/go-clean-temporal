package domain

import (
	"net/http"

	"github.com/lengocson131002/go-clean-core/errors"
)

var (
	// DOMAIN CUSTOM ERROR
	ErrorAccountNotFound = &errors.DomainError{
		Status:  http.StatusBadRequest, // http mapping
		Code:    "100",
		Message: "User not found",
	}
)
