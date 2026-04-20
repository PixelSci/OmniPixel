package apperr

import "net/http"

// Predefined errors act as roots. Treat them as immutable: use WithDetails /
// WithCause / WithMessage to derive per-call copies.
var (
	// 10xxx — common / system
	ErrUnknown         = New(CodeUnknown, http.StatusInternalServerError, "unknown error")
	ErrInternal        = New(CodeInternal, http.StatusInternalServerError, "internal server error")
	ErrInvalidParam    = New(CodeInvalidParam, http.StatusBadRequest, "invalid parameter")
	ErrValidation      = New(CodeValidation, http.StatusBadRequest, "validation failed")
	ErrNotFound        = New(CodeNotFound, http.StatusNotFound, "resource not found")
	ErrConflict        = New(CodeConflict, http.StatusConflict, "resource conflict")
	ErrTooManyRequests = New(CodeTooManyRequests, http.StatusTooManyRequests, "too many requests")
	ErrTimeout         = New(CodeTimeout, http.StatusGatewayTimeout, "request timeout")
	ErrUnavailable     = New(CodeUnavailable, http.StatusServiceUnavailable, "service unavailable")

	// 11xxx — auth
	ErrUnauthorized = New(CodeUnauthorized, http.StatusUnauthorized, "unauthorized")
	ErrTokenExpired = New(CodeTokenExpired, http.StatusUnauthorized, "token expired")
	ErrTokenInvalid = New(CodeTokenInvalid, http.StatusUnauthorized, "token invalid")
	ErrForbidden    = New(CodeForbidden, http.StatusForbidden, "forbidden")

	// 12xxx — user
	ErrUserNotFound  = New(CodeUserNotFound, http.StatusNotFound, "user not found")
	ErrUserExists    = New(CodeUserExists, http.StatusConflict, "user already exists")
	ErrPasswordWrong = New(CodePasswordWrong, http.StatusUnauthorized, "incorrect password")
	ErrUserDisabled  = New(CodeUserDisabled, http.StatusForbidden, "user disabled")
)
