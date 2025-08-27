package rest_err

import (
	"l03/internal/internal_error"
	"net/http"
)

type RestError struct {
	Message string   `json:"message"`
	Code    int      `json:"code"`
	Err     string   `json:"err"`
	Causes  []Causes `json:"causes"`
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestError) Error() string {
	return r.Message
}

func NewRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     "Bad Request",
		Causes:  nil,
	}
}

func NewInternalServerError(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     "internal_server",
		Causes:  nil,
	}
}

func NewNotFoundError(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusNotFound,
		Err:     "not_found",
		Causes:  nil,
	}
}

func NewBadRequestError(message string, causes ...Causes) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     "bad_request",
		Causes:  causes,
	}
}

func ConvertError(ir *internal_error.InternalError) *RestError {
	switch ir.Err {
	case "bad_request":
		return NewBadRequestError(ir.Error())
	case "not_found":
		return NewNotFoundError(ir.Error())
	default:
		return NewInternalServerError(ir.Error())
	}
}
