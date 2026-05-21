package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewConversationRoutes(router fiber.Router, conversationController *controller.ConversationController) {
	router.Post("/conversation", func() {})
	router.Get("/conversations", conversationController.ListConversations)
	router.Get("/conversations/:conversation_id", conversationController.GetConversation)
}
