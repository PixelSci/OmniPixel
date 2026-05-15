package controller

import (
	"omni-pixel/usecase"

	"github.com/gofiber/fiber/v3"
)

type HealthController struct {
	healthUseCase *usecase.HealthUseCase
}

func NewHealthController(healthUseCase *usecase.HealthUseCase) *HealthController {
	return &HealthController{healthUseCase}
}

func (healthController *HealthController) Health(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": "healthy"})
}
