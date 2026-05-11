package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
	"omni-pixel/api/middleware"
)

func NewSessionRoute(router fiber.Router, sessionController *controller.SessionController, accessTokenSecret string) {
	router.Get("/sessions", middleware.JWTAuthMiddleware(accessTokenSecret), sessionController.ListSessions)
}
