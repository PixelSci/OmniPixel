package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
	"omni-pixel/api/middleware"
)

func NewSessionRoute(router fiber.Router, sessionController *controller.SessionController, accessTokenSecret string) {
	router.Get("/sessions", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.ListSessions)
	router.Post("/sessions", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.CreateSession)
	router.Get("/sessions/:id", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.GetSession)
	router.Post("/sessions/:id/chat", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.SendSessionPrompt)
	router.Post("/sessions/:id/content", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.SaveSessionChatContent)
	router.Delete("/sessions/:id", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.DeleteSession)
}
