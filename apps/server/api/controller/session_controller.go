package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"omni-pixel/domain"
	"omni-pixel/usecase"
)

type SessionController struct {
	sessionUsecase *usecase.SessionUsecase
}

func NewSessionController(sessionUsecase *usecase.SessionUsecase) *SessionController {
	return &SessionController{sessionUsecase: sessionUsecase}
}

func (controller *SessionController) ListSessions(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
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

func (controller *SessionController) DeleteSession(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	if err := controller.sessionUsecase.DeleteSession(c.Params("id"), userID); errors.Is(err, domain.ErrSessionNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "session not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete session",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func userIDFromLocals(c fiber.Ctx) (string, error) {
	userID, ok := fiber.Locals(c, "user_id").(string)
	if !ok || userID == "" {
		return "", domain.ErrInvalidCredentials
	}

	return userID, nil
}
