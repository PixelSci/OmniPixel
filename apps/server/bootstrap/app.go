package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/api/controller"
	"omni-pixel/api/routes"
	"omni-pixel/internal"
	"omni-pixel/repository"
	"omni-pixel/usecase"
)

func NewApp(env *Env, db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "OmniPixel Server",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	timeout := time.Duration(env.ContextTimeout) * time.Second
	userRepository := repository.NewUserRepository(db, timeout)
	sessionRepository := repository.NewSessionRepository(db, timeout)
	chatCompletionClient := internal.NewChatCompletionClient(timeout)
	signinUsecase := usecase.NewSigninUsecase(
		userRepository,
		env.AccessTokenSecret,
		env.AccessTokenExpiryHour,
	)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepository, chatCompletionClient)
	signinController := controller.NewSigninController(signinUsecase)
	sessionController := controller.NewSessionController(sessionUsecase)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	routes.NewSigninRoute(v1, signinController)
	routes.NewSessionRoute(v1, sessionController, env.AccessTokenSecret)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return app
}
