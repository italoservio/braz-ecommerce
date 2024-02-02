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
)

func Http(content ...string) *HTTPException {
	if len(content) == 0 {
		return &HTTPException{ErrorMessage: "Invalid error code"}
	}

	return errorCodeToStruct(content[0])
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
		response.ErrorMessage = "Failed to communicated with database"
	case CodeValidationFailed:
		response.StatusMessage = "Bad Request"
		response.StatusCode = http.StatusBadRequest
		response.ErrorMessage = "Invalid input for one or more required attributes"
	}

	return &response
}
