package controller

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/usecase"
)

type SessionController struct {
	sessionUsecase *usecase.SessionUsecase
}

func NewSessionController(sessionUsecase *usecase.SessionUsecase) *SessionController {
	return &SessionController{sessionUsecase: sessionUsecase}
}

func (controller *SessionController) ListSessions(c fiber.Ctx) error {
	userID, ok := fiber.Locals(c, "user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	response, err := controller.sessionUsecase.ListSessions(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to list sessions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
