package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"omni-pixel/apperr"
	"omni-pixel/auth"
	"omni-pixel/config"
	"omni-pixel/db"
	"omni-pixel/handler"
	"omni-pixel/repository"
	"omni-pixel/router"
	"omni-pixel/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.NewPool(ctx, db.Config{
		DSN:             cfg.DB.DSN,
		MaxConns:        cfg.DB.MaxConns,
		MinConns:        cfg.DB.MinConns,
		MaxConnLifetime: cfg.DB.MaxConnLifetime,
		MaxConnIdleTime: cfg.DB.MaxConnIdleTime,
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if cfg.DB.AutoMigrate {
		if err := db.Migrate(cfg.DB.DSN); err != nil {
			log.Fatalf("db migrate: %v", err)
		}
		log.Println("db migrations applied")
	}

	userRepo := repository.NewPostgresUserRepo(pool)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	tokenIssuer := auth.NewTokenIssuer(cfg.Auth.JWTSecret, cfg.Auth.JWTTTL, cfg.Auth.JWTIssuer)
	passwordProvider := auth.NewPasswordProvider(userRepo)
	authSvc := auth.NewService(userRepo, tokenIssuer, passwordProvider)
	// When GitHub / Google providers land, register them here:
	//   authSvc.RegisterOAuth(github.New(cfg.Auth.GitHub))
	//   authSvc.RegisterOAuth(google.New(cfg.Auth.Google))
	authHandler := handler.NewAuthHandler(authSvc)

	app := fiber.New(fiber.Config{
		AppName:      "omni-pixel",
		ErrorHandler: apperr.FiberErrorHandler,
	})
	app.Use(
		recover.New(),
		requestid.New(),
		logger.New(logger.Config{
			Format: "[${time}] ${locals:requestid} ${status} ${method} ${path} (${latency})\n",
		}),
	)

	router.Register(app, userHandler, authHandler)

	go func() {
		if err := app.Listen(cfg.Addr); err != nil {
			log.Printf("listen error: %v", err)
		}
	}()
	log.Printf("server listening on %s (env=%s)", cfg.Addr, cfg.Env)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
