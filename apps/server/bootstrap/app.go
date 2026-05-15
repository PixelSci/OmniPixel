package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/api/controller"
	"omni-pixel/api/middleware"
	"omni-pixel/api/routes"
	"omni-pixel/repository"
	"omni-pixel/usecase"
)

func NewApp(env *Env, db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "OmniPixel Server",
	})

	app.Use(middleware.CORSMiddleware())

	timeout := time.Duration(env.ContextTimeout) * time.Second
	userRepository := repository.NewUserRepository(db, timeout)
	sessionRepository := repository.NewSessionRepository(db, timeout)

	healthUseCase := usecase.NewHealthUseCase()
	accountUseCase := usecase.NewAccountUseCase(
		userRepository,
		env.AccessTokenSecret,
		env.AccessTokenExpiryHour,
	)
	sessionUseCase := usecase.NewSessionUseCase(sessionRepository)

	healthController := controller.NewHealthController(healthUseCase)
	accountController := controller.NewAccountController(accountUseCase)
	sessionController := controller.NewSessionController(sessionUseCase)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.NewHealthRoute(v1, healthController)
	routes.NewAccountRoutes(v1, accountController)

	auth := v1.Group("", middleware.JWTAuthMiddleware(env.AccessTokenSecret))
	routes.NewSessionRoutes(auth, sessionController)

	return app
}
