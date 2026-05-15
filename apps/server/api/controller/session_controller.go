package controller

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/usecase"
)

type SessionController struct {
	sessionUseCase *usecase.SessionUseCase
}

func NewSessionController(sessionUseCase *usecase.SessionUseCase) *SessionController {
	return &SessionController{sessionUseCase}
}

func (controller *SessionController) ListSessions(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": "success"})
}
