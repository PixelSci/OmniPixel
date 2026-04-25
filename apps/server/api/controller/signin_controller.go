package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"omni-pixel/domain"
	"omni-pixel/usecase"
)

type SigninController struct {
	signinUsecase *usecase.SigninUsecase
}

func NewSigninController(signinUsecase *usecase.SigninUsecase) *SigninController {
	return &SigninController{signinUsecase: signinUsecase}
}

func (controller *SigninController) Signin(c fiber.Ctx) error {
	var request domain.SigninRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	response, err := controller.signinUsecase.Signin(request)
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
