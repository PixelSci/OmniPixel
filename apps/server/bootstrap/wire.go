package bootstrap

import (
	"log"
	"time"

	"omni-pixel/api/controller"
	"omni-pixel/repository"
	"omni-pixel/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Providers struct {
	ENV               *Env
	HealthController  *controller.HealthController
	AccountController *controller.AccountController
	SessionController *controller.SessionController
}

func NewProviders() *Providers {
	env := newEnv()
	db := newDB(env)
	defer db.Close()

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

	return &Providers{
		ENV:               env,
		HealthController:  controller.NewHealthController(healthUseCase),
		AccountController: controller.NewAccountController(accountUseCase),
		SessionController: controller.NewSessionController(sessionUseCase),
	}
}

func newEnv() *Env {
	env, err := NewEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	return env
}

func newDB(env *Env) *pgxpool.Pool {
	db, err := NewPostgresPool(env)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
