package routes

import (
	"github.com/gofiber/fiber/v3"

	"omni-pixel/api/controller"
)

func NewAccountRoutes(router fiber.Router, accountController *controller.AccountController) {
	router.Post("/signin", accountController.Signin)
	router.Post("/signup", accountController.Signup)
}
