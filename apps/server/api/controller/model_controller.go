package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"omni-pixel/domain"
	"omni-pixel/internal/response"
	"omni-pixel/usecase"
)

type ModelController struct {
	modelUseCase *usecase.ModelUseCase
}

func NewModelController(modelUseCase *usecase.ModelUseCase) *ModelController {
	return &ModelController{modelUseCase}
}

func (controller *ModelController) ListModels(c fiber.Ctx) error {
	filter := domain.ModelFilter{}

	if providerID := c.Query("provider_id"); providerID != "" {
		id, err := uuid.Parse(providerID)
		if err != nil {
			return response.Write(c, response.ErrInvalidRequest)
		}
		filter.ProviderID = &id
	}

	if c.Query("include_disabled") == "true" {
		filter.IncludeDisabled = true
	}

	models, err := controller.modelUseCase.ListModels(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"models": models})
}
