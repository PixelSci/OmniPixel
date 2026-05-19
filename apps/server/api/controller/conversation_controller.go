package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

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
