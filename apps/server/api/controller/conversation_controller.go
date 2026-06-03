package controller

import (
	"bufio"
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

func (controller *ConversationController) DeleteConversation(c fiber.Ctx) error {
	conversationID, err := uuid.Parse(c.Params("conversation_id"))
	if err != nil {
		return response.Write(c, response.ErrInvalidConvID)
	}

	userID := fiber.Locals[uuid.UUID](c, "user_id")
	if err := controller.conversationUseCase.DeleteConversation(userID, conversationID); err != nil {
		return response.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type sseStreamWriter struct {
	w *bufio.Writer
}

func (s *sseStreamWriter) WriteToken(token string) error {
	data, _ := json.Marshal(map[string]string{"token": token})
	_, err := fmt.Fprintf(s.w, "data: %s\n\n", data)
	if err != nil {
		return err
	}
	return s.w.Flush()
}

func (s *sseStreamWriter) WriteDone(conversationID, messageID uuid.UUID) error {
	data, _ := json.Marshal(map[string]interface{}{
		"done":            true,
		"conversation_id": conversationID.String(),
		"message_id":      messageID.String(),
	})
	_, err := fmt.Fprintf(s.w, "data: %s\n\n", data)
	if err != nil {
		return err
	}
	return s.w.Flush()
}

func (controller *ConversationController) Chat(c fiber.Ctx) error {
	var request domain.ChatRequest
	if err := c.Bind().Body(&request); err != nil {
		return response.Write(c, response.ErrInvalidRequest)
	}

	userID := fiber.Locals[uuid.UUID](c, "user_id")

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	return c.Status(fiber.StatusOK).SendStreamWriter(func(w *bufio.Writer) {
		writer := &sseStreamWriter{w: w}
		if err := controller.conversationUseCase.Chat(userID, request, writer); err != nil {
			data, _ := json.Marshal(map[string]string{"error": err.Error()})
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.Flush()
		}
	})
}
