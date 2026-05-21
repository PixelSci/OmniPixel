package bootstrap

import (
	"log"

	"omni-pixel/api/controller"
	"omni-pixel/internal/ai"
	"omni-pixel/repository"
	"omni-pixel/usecase"

	"gorm.io/gorm"
)

type Providers struct {
	ENV                    *Env
	HealthController       *controller.HealthController
	AccountController      *controller.AccountController
	ConversationController *controller.ConversationController
}

func NewProviders() *Providers {
	env := newEnv()
	db := newDB(env)

	userRepository := repository.NewUserRepository(db)
	conversationRepository := repository.NewConversationRepository(db)

	healthUseCase := usecase.NewHealthUseCase()
	accountUseCase := usecase.NewAccountUseCase(
		userRepository,
		env.AccessTokenSecret,
		env.AccessTokenExpiryHour,
	)
	aiProvider := ai.NewOpenAIProvider(env.AIBaseURL, env.AIAPIKey)
	conversationUseCase := usecase.NewConversationUseCase(conversationRepository, aiProvider)

	return &Providers{
		ENV:                    env,
		HealthController:       controller.NewHealthController(healthUseCase),
		AccountController:      controller.NewAccountController(accountUseCase),
		ConversationController: controller.NewConversationController(conversationUseCase),
	}
}

func newEnv() *Env {
	env, err := NewEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	return env
}

func newDB(env *Env) *gorm.DB {
	db, err := NewPostgresDB(env)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
