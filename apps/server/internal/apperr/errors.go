// Package apperr provides a structured, JSON-friendly error type for the API server.
//
// AppError carries a business code, HTTP status, human message and optional
// details/cause. Use the predefined variables in predefined.go as roots and
// chain With* helpers to attach context without mutating the original.
package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is the canonical error returned by service and handler layers.
//
// Fields are exported for JSON marshaling; do not mutate them directly on
// predefined values — use the With* methods, which return copies.
type AppError struct {
	Code     int            `json:"code"`
	HTTPCode int            `json:"-"`
	Message  string         `json:"message"`
	Details  map[string]any `json:"details,omitempty"`
	cause    error
}

// New builds a fresh AppError.
func New(code, httpCode int, message string) *AppError {
	return &AppError{Code: code, HTTPCode: httpCode, Message: message}
}

// Newf is New with fmt.Sprintf-style formatting.
func Newf(code, httpCode int, format string, args ...any) *AppError {
	return &AppError{Code: code, HTTPCode: httpCode, Message: fmt.Sprintf(format, args...)}
}

// Wrap attaches an AppError identity on top of an underlying cause.
// Returns nil when err is nil so callers can chain without nil checks.
func Wrap(err error, code, httpCode int, message string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{Code: code, HTTPCode: httpCode, Message: message, cause: err}
}

// From normalizes any error into an *AppError. Already-typed errors pass
// through; anything else becomes an opaque internal error wrapping the cause.
func From(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return &AppError{
		Code:     CodeInternal,
		HTTPCode: http.StatusInternalServerError,
		Message:  err.Error(),
		cause:    err,
	}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As traversal.
func (e *AppError) Unwrap() error { return e.cause }

// Is matches by Code so errors.Is(err, apperr.ErrUserNotFound) still works
// after the error has been wrapped with additional context.
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// WithDetails returns a copy with the given details merged in. Keys in the
// argument override existing keys; the original map is never mutated.
func (e *AppError) WithDetails(details map[string]any) *AppError {
	out := *e
	merged := make(map[string]any, len(e.Details)+len(details))
	for k, v := range e.Details {
		merged[k] = v
	}
	for k, v := range details {
		merged[k] = v
	}
	out.Details = merged
	return &out
}

// WithDetail is WithDetails for a single key/value pair.
func (e *AppError) WithDetail(key string, value any) *AppError {
	return e.WithDetails(map[string]any{key: value})
}

// WithCause returns a copy with the underlying cause set/replaced.
func (e *AppError) WithCause(cause error) *AppError {
	out := *e
	out.cause = cause
	return &out
}

// WithMessage returns a copy with the user-facing message replaced.
func (e *AppError) WithMessage(message string) *AppError {
	out := *e
	out.Message = message
	return &out
}

// WithMessagef is WithMessage with fmt.Sprintf-style formatting.
func (e *AppError) WithMessagef(format string, args ...any) *AppError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}
