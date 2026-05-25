package apperr

import "net/http"

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrSeatTaken = &AppError{
		Code:       "SEAT_TAKEN",
		Message:    "seat is already occupied",
		HTTPStatus: http.StatusConflict,
	}
	ErrSeatNotFound = &AppError{
		Code:       "SEAT_NOT_FOUND",
		Message:    "seat not found",
		HTTPStatus: http.StatusNotFound,
	}
	ErrShowNotFound = &AppError{
		Code:       "SHOW_NOT_FOUND",
		Message:    "show not found",
		HTTPStatus: http.StatusNotFound,
	}
	ErrInvalidRequestBody = &AppError{
		Code:       "INVALID_REQUEST_BODY",
		Message:    "invalid request body",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrInvalidSeatCount = &AppError{
		Code:       "INVALID_SEAT_COUNT",
		Message:    "seat count must be greater than 0",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "missing X-User-ID header",
		HTTPStatus: http.StatusUnauthorized,
	}
)

func ErrInvalidID(name string) *AppError {
	return &AppError{
		Code:       "INVALID_ID",
		Message:    "invalid " + name + " id",
		HTTPStatus: http.StatusBadRequest,
	}
}
