package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewModelRoutes(router fiber.Router, modelController *controller.ModelController) {
	router.Get("/models", modelController.ListModels)
}
