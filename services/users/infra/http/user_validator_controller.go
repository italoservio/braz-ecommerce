package http

import (
	"errors"

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

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
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
	myValidator := &XValidator{
		validator: validate,
	}

	if errs := myValidator.Validate(payloadValidate); len(errs) > 0 && errs[0].Error {
		return errors.New(exception.CodeValidationFailed)
	}

	return nil
}
