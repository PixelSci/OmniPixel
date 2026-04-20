package handler

import "github.com/gofiber/fiber/v2"

// SuccessResponse is the JSON envelope used for successful responses; it mirrors
// apperr.Response so clients see a consistent {code, message, data|details} shape.
type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ok(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(SuccessResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Success(c *fiber.Ctx, data any) error { return ok(c, fiber.StatusOK, data) }
func Created(c *fiber.Ctx, data any) error { return ok(c, fiber.StatusCreated, data) }
