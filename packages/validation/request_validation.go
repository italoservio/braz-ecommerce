package validation

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
)

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

func ValidateRequest(c *fiber.Ctx, payload any) error {
	correlationId := c.Locals(string(logger.CorrelationId))
	validate := validator.New()
	validationErrors := []ErrorResponse{}
	errs := validate.Struct(payload)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			slog.Error(fmt.Sprintf("%s %s", correlationId, err.Error()))

			validationErrors = append(validationErrors, ErrorResponse{
				Error:       true,
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
			})
		}
	}
	if len(validationErrors) > 0 {
		return errors.New(exception.CodeValidationFailed)
	}
	return nil
}
