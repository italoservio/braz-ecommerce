package http

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}
)

func (v XValidator) Validate(data any, validate *validator.Validate) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {

			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func ValidationRequest(c *fiber.Ctx, payloadValidate interface{}) error {
	validate := validator.New()

	myValidator := &XValidator{
		validator: validate,
	}

	if errs := myValidator.Validate(payloadValidate, validate); len(errs) > 0 && errs[0].Error {
		slog.Error("missing information in body")
		return errors.New(exception.CodeInternal)
	}

	return nil
}
