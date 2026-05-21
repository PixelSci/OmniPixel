package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"omni-pixel/domain"
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

type sseAdapter struct {
	c fiber.Ctx
}

func (w *sseAdapter) WriteToken(token string) error {
	data, _ := json.Marshal(map[string]string{"token": token})
	_, err := fmt.Fprintf(w.c, "data: %s\n\n", data)
	return err
}

func (w *sseAdapter) WriteDone(conversationID, messageID uuid.UUID) error {
	data, _ := json.Marshal(map[string]interface{}{
		"done":            true,
		"conversation_id": conversationID.String(),
		"message_id":      messageID.String(),
	})
	_, err := fmt.Fprintf(w.c, "data: %s\n\n", data)
	return err
}

func (controller *ConversationController) Chat(c fiber.Ctx) error {
	var request domain.ChatRequest
	if err := c.Bind().Body(&request); err != nil {
		return response.Write(c, response.ErrInvalidRequest)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	userID := fiber.Locals[uuid.UUID](c, "user_id")
	writer := &sseAdapter{c: c}

	if err := controller.conversationUseCase.Chat(userID, request, writer); err != nil {
		return response.DomainError(c, err)
	}

	return nil
}
