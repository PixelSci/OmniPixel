package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"omni-pixel/domain"
	"omni-pixel/usecase"
)

type ConversationController struct {
	conversationUseCase *usecase.ConversationUseCase
}

func NewConversationController(conversationUseCase *usecase.ConversationUseCase) *ConversationController {
	return &ConversationController{conversationUseCase}
}

func (controller *ConversationController) ListConversations(c fiber.Ctx) error {
	userId := fiber.Locals[uuid.UUID](c, "user_id")
	conversations, err := controller.conversationUseCase.ListConversations(userId)
	if err != nil {
		return err
	}

	return c.JSON(conversations)
}

func (controller *ConversationController) GetConversation(c fiber.Ctx) error {
	conversationID, err := uuid.Parse(c.Params("conversation_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid conversation id",
		})
	}

	userID := fiber.Locals[uuid.UUID](c, "user_id")
	response, err := controller.conversationUseCase.GetConversation(userID, conversationID)
	if errors.Is(err, domain.ErrInvalidConversationID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid conversation id",
		})
	}
	if errors.Is(err, domain.ErrConversationNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "conversation not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get conversation",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
