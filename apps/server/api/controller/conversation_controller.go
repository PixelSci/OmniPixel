package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"omni-pixel/internal/response"
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
		return response.Write(c, response.ErrInvalidConvID)
	}

	userID := fiber.Locals[uuid.UUID](c, "user_id")
	result, err := controller.conversationUseCase.GetConversation(userID, conversationID)
	if err != nil {
		return response.DomainError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
