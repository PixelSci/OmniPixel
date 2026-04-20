package router

import (
	"github.com/gofiber/fiber/v2"

	"omni-pixel/handler"
)

func Register(app *fiber.App, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	// Future OAuth routes (github / google) will land here, e.g.:
	//   auth.Get("/:provider/url",      authHandler.OAuthRedirect)
	//   auth.Get("/:provider/callback", authHandler.OAuthCallback)

	users := api.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
}
