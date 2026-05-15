package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"omni-pixel/domain"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	response, err := controller.accountUseCase.Signin(request)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid email or password",
		})
	}
	if errors.Is(err, domain.ErrUserInactive) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "user is not active",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to sign in",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (controller *AccountController) Signup(c fiber.Ctx) error {
	var request domain.SigninRequest

	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	return nil
}
