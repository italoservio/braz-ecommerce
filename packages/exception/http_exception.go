package exception

import (
	"net/http"
)

type HTTPException struct {
	Ok            bool   `json:"ok"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	ErrorCode     string `json:"error_code"`
	ErrorMessage  string `json:"error_message"`
}

const (
	CodeNotFound         = "ENOTFOUND"
	CodeDatabaseFailed   = "EDBFAILURE"
	CodeValidationFailed = "EVALIDATION"
	CodeInternal         = "EINTERNAL"
	CodePermission       = "EPERMISSION"
)

func Http(code string) *HTTPException {
	codes := map[string]bool{
		(CodeNotFound):         true,
		(CodeDatabaseFailed):   true,
		(CodeValidationFailed): true,
		(CodeInternal):         true,
		(CodePermission):       true,
	}

	if !codes[code] {
		return &HTTPException{
			ErrorCode:     "EINVALID",
			ErrorMessage:  "Invalid error code",
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "Internal Server Error",
		}
	}

	return errorCodeToStruct(code)
}

func (h *HTTPException) Error() string {
	return h.ErrorMessage
}

func errorCodeToStruct(code string) *HTTPException {
	response := HTTPException{Ok: false, ErrorCode: code}

	switch code {
	case CodeNotFound:
		response.StatusMessage = "Not Found"
		response.StatusCode = http.StatusNotFound
		response.ErrorMessage = "Entity not found"
	case CodeDatabaseFailed:
		response.StatusMessage = "Internal Server Error"
		response.StatusCode = http.StatusInternalServerError
		response.ErrorMessage = "Failed to communicate with database"
	case CodeValidationFailed:
		response.StatusMessage = "Bad Request"
		response.StatusCode = http.StatusBadRequest
		response.ErrorMessage = "Invalid input for one or more required attributes"
	case CodeInternal:
		response.StatusMessage = "Internal Server Error"
		response.StatusCode = http.StatusInternalServerError
		response.ErrorMessage = "An expected error occurred and the server could not deal with it"
	case CodePermission:
		response.StatusMessage = "Unauthorized"
		response.StatusCode = http.StatusUnauthorized
		response.ErrorMessage = "User not allowed to perform this action"
	}

	return &response
}
