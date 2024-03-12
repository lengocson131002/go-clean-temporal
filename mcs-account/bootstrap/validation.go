package bootstrap

import (
	"github.com/lengocson131002/go-clean-core/validation"
	"github.com/lengocson131002/go-clean-core/validation/goplayaround"
)

func GetValidator() validation.Validator {
	return goplayaround.NewGpValidator()
}
