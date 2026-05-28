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
	ProviderController     *controller.ProviderController
	ModelController        *controller.ModelController
}

func NewProviders() *Providers {
	env := newEnv()
	db := newDB(env)

	userRepository := repository.NewUserRepository(db)
	conversationRepository := repository.NewConversationRepository(db)
	providerRepository := repository.NewProviderRepository(db)
	modelRepository := repository.NewModelRepository(db)

	healthUseCase := usecase.NewHealthUseCase()
	accountUseCase := usecase.NewAccountUseCase(
		userRepository,
		env.AccessTokenSecret,
		env.AccessTokenExpiryHour,
	)
	aiProvider := ai.NewOpenAIProvider(env.AIBaseURL, env.AIAPIKey)
	conversationUseCase := usecase.NewConversationUseCase(conversationRepository, aiProvider)
	providerUseCase := usecase.NewProviderUseCase(providerRepository)
	modelUseCase := usecase.NewModelUseCase(modelRepository)

	return &Providers{
		ENV:                    env,
		HealthController:       controller.NewHealthController(healthUseCase),
		AccountController:      controller.NewAccountController(accountUseCase),
		ConversationController: controller.NewConversationController(conversationUseCase),
		ProviderController:     controller.NewProviderController(providerUseCase),
		ModelController:        controller.NewModelController(modelUseCase),
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
