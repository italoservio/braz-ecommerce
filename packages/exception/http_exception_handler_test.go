package exception

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestHttpExceptionHandler(t *testing.T) {
	t.Run("should return http exception when code is correct", func(t *testing.T) {
		fbr := fiber.New()
		ctx := fbr.AcquireCtx(&fasthttp.RequestCtx{})
		err := HttpExceptionHandler(ctx, errors.New(CodeValidationFailed))
		assert.Nil(t, err, "should not return error")
	})
}
