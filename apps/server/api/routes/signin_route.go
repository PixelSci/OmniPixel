package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewSigninRoute(router fiber.Router, signinController *controller.SigninController) {
	router.Post("/signin", signinController.Signin)
}
