package controller

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/domain"
	"omni-pixel/internal/response"
	"omni-pixel/usecase"
)

type AccountController struct {
	accountUseCase *usecase.AccountUseCase
}

func NewAccountController(accountUseCase *usecase.AccountUseCase) *AccountController {
	return &AccountController{accountUseCase}
}

func (controller *AccountController) Signin(c fiber.Ctx) error {
	var request domain.SigninRequest
	if err := c.Bind().Body(&request); err != nil {
		return response.Write(c, response.ErrInvalidRequest)
	}

	result, err := controller.accountUseCase.Signin(request)
	if err != nil {
		return response.DomainError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *AccountController) Signup(c fiber.Ctx) error {
	var request domain.SignupRequest
	if err := c.Bind().Body(&request); err != nil {
		return response.Write(c, response.ErrInvalidRequest)
	}

	result, err := controller.accountUseCase.Signup(request)
	if err != nil {
		return response.DomainError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
