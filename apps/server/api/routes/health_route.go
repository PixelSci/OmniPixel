package routes

import (
	"omni-pixel/api/controller"

	"github.com/gofiber/fiber/v3"
)

func NewHealthRoute(router fiber.Router, healthController *controller.HealthController) {
	router.Get("/health", healthController.Health)
}
