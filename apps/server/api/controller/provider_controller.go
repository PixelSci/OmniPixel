package controller

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/usecase"
)

type ProviderController struct {
	providerUseCase *usecase.ProviderUseCase
}

func NewProviderController(providerUseCase *usecase.ProviderUseCase) *ProviderController {
	return &ProviderController{providerUseCase}
}

func (controller *ProviderController) ListProviders(c fiber.Ctx) error {
	providers, err := controller.providerUseCase.ListProviders()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"providers": providers})
}
