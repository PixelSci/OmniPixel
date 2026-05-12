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

func (controller *SessionController) GetSession(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	response, err := controller.sessionUsecase.GetSession(c.Params("id"), userID)
	if errors.Is(err, domain.ErrSessionNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "session not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (controller *SessionController) CreateSession(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	var request domain.CreateSessionRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	request.UserID = userID

	response, err := controller.sessionUsecase.CreateSession(request)
	if errors.Is(err, domain.ErrInvalidSessionID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid session id",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create session",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (controller *SessionController) SaveSessionChatContent(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	var request domain.SaveSessionChatContentRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	request.SessionID = c.Params("id")
	request.UserID = userID

	if err := controller.sessionUsecase.SaveSessionChatContent(request); errors.Is(err, domain.ErrSessionNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "session not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to save session chat content",
		})
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (controller *SessionController) SendSessionPrompt(c fiber.Ctx) error {
	userID, err := userIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	var request domain.SendSessionPromptRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	if request.SessionID == "" {
		request.SessionID = c.Params("id")
	}
	request.UserID = userID

	response, err := controller.sessionUsecase.SendSessionPrompt(request)
	if errors.Is(err, domain.ErrInvalidPrompt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "prompt is required",
		})
	}
	if errors.Is(err, domain.ErrInvalidSessionID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid session id",
		})
	}
	if errors.Is(err, domain.ErrInvalidAIConfig) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ai config",
		})
	}
	if errors.Is(err, domain.ErrUnsupportedAIProvider) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "unsupported ai provider",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to send session prompt",
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

	c.Status(fiber.StatusNoContent)
	return nil
}

func userIDFromLocals(c fiber.Ctx) (string, error) {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return "", domain.ErrInvalidCredentials
	}

	return userID, nil
}
