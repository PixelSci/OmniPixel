package bootstrap

import (
	"log"
	"omni-pixel/api/middleware"
	"omni-pixel/api/routes"

	"github.com/gofiber/fiber/v3"
)

type BootStrap struct {
	app *fiber.App
	p   *Providers
}

func NewApp() *BootStrap {
	app := fiber.New(fiber.Config{
		AppName: "OmniPixel Server",
	})

	app.Use(middleware.CORSMiddleware())

	providers := NewProviders()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.NewHealthRoute(v1, providers.HealthController)
	routes.NewAccountRoutes(v1, providers.AccountController)

	auth := v1.Group("", middleware.JWTAuthMiddleware(providers.ENV.AccessTokenSecret))
	routes.NewConversationRoutes(auth, providers.ConversationController)

	return &BootStrap{app, providers}
}

func (b *BootStrap) Start() {
	if err := b.app.Listen(b.p.ENV.ServerAddress); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
