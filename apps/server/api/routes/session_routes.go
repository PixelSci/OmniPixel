package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewSessionRoutes(router fiber.Router, sessionController *controller.SessionController) {
	router.Get("/sessions/list", sessionController.ListSessions)
}
