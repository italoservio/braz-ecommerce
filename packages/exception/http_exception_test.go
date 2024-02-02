package exception

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestException_Http(t *testing.T) {
	t.Run("should parse code to http exception structure", func(t *testing.T) {
		structure := Http(CodeNotFound)

		assert.Equal(t, structure.Ok, false)
		assert.Equal(t, structure.ErrorCode, CodeNotFound)
		assert.Equal(t, structure.ErrorMessage, "Entity not found")
		assert.Equal(t, structure.StatusCode, 404)
		assert.Equal(t, structure.StatusMessage, "Not Found")
	})

	t.Run("should parse invalid code to an empty http exception structure", func(t *testing.T) {
		structure := Http("")

		fmt.Printf("%+v", structure)

		assert.Equal(t, false, structure.Ok)
		assert.Equal(t, "", structure.ErrorCode)
		assert.Equal(t, "Invalid error code", structure.ErrorMessage)
		assert.Equal(t, 0, structure.StatusCode)
		assert.Equal(t, "", structure.StatusMessage)
	})
}

func TestException_errorCodeToStruct(t *testing.T) {
	t.Run("should parse error code ENOTFOUND", func(t *testing.T) {
		structure := errorCodeToStruct(CodeNotFound)

		assert.Equal(t, structure.StatusCode, 404)
		assert.Equal(t, structure.ErrorMessage, "Entity not found")
	})

	t.Run("should parse error code EDBFAILURE", func(t *testing.T) {
		structure := errorCodeToStruct(CodeDatabaseFailed)

		assert.Equal(t, structure.StatusCode, 500)
		assert.Equal(t, structure.ErrorMessage, "Failed to communicate with database")
	})

	t.Run("should parse error code EDBFAILURE", func(t *testing.T) {
		structure := errorCodeToStruct(CodeValidationFailed)

		assert.Equal(t, structure.StatusCode, 400)
		assert.Equal(t, structure.ErrorMessage, "Invalid input for one or more required attributes")
	})
}
