package create_user

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{ createUserService CreateUserServiceInterface }

func NewUserController(createUserService CreateUserServiceInterface) *UserController {
	return &UserController{createUserService}
}

type (
	User struct {
		Name string `validate:"required,min=5,max=20"` // Required field, min 5 char long max 20
		Age  int    `validate:"required,teener"`       // Required field, and client needs to implement our 'teener' tag format which we'll see later
	}

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
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func (cr *UserController) CreateUser(c *fiber.Ctx) error {
	payloadValidate := DTOCreateUserReq{}

	if err := c.BodyParser(&payloadValidate); err != nil {
		return err
	}
	myValidator := &XValidator{
		validator: validate,
	}
	myValidator.validator.RegisterValidation("teener", func(fl validator.FieldLevel) bool {
		// User.Age needs to fit our needs, 12-18 years old.
		return fl.Field().Int() >= 12 && fl.Field().Int() <= 18
	})

	if errs := myValidator.Validate(payloadValidate); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	// Logic, validated with success
	return c.SendString("Hello, World!")
	// if err := c.BodyParser(&payload); err != nil {
	// 	return err
	// }

	// return c.JSON(payload)
}
