package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/api/controller"
	"omni-pixel/api/routes"
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
	routes.NewSessionRoutes(v1, sessionController)

	return app
}
