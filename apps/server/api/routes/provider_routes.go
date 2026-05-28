package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewProviderRoutes(router fiber.Router, providerController *controller.ProviderController) {
	router.Get("/providers", providerController.ListProviders)
}
