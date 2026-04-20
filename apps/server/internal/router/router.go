package router

import (
	"github.com/gofiber/fiber/v2"

	"omni-pixel/internal/handler"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")

	users := api.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
}
