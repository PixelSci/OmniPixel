package apperr

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Response is the unified JSON shape returned to clients on error.
type Response struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Render writes err as a standardized JSON response. Non-AppError values are
// normalized via From (becoming CodeInternal / HTTP 500).
func Render(c *fiber.Ctx, err error) error {
	appErr := From(err)
	return c.Status(appErr.HTTPCode).JSON(Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	})
}

// FiberErrorHandler is intended for fiber.Config{ ErrorHandler: ... }. It
// translates fiber's own *fiber.Error (e.g. 404 / 405 from the router) into
// the standardized Response shape and delegates everything else to Render.
func FiberErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(Response{
			Code:    mapHTTPToCode(fiberErr.Code),
			Message: fiberErr.Message,
		})
	}
	return Render(c, err)
}

func mapHTTPToCode(httpStatus int) int {
	switch httpStatus {
	case http.StatusBadRequest:
		return CodeInvalidParam
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
	case http.StatusGatewayTimeout:
		return CodeTimeout
	case http.StatusServiceUnavailable:
		return CodeUnavailable
	default:
		return CodeInternal
	}
}
