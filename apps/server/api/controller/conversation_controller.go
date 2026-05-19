package controller

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/usecase"
)

type ConversationController struct {
	conversationUseCase *usecase.ConversationUseCase
}

func NewConversationController(conversationUseCase *usecase.ConversationUseCase) *ConversationController {
	return &ConversationController{conversationUseCase}
}

func (controller *ConversationController) ListConversations(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": "success"})
}
