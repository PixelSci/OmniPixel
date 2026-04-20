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

	"omni-pixel/internal/apperr"
	"omni-pixel/internal/config"
	"omni-pixel/internal/handler"
	"omni-pixel/internal/repository"
	"omni-pixel/internal/router"
	"omni-pixel/internal/service"
)

func main() {
	cfg := config.Load()

	userRepo := repository.NewInMemoryUserRepo()
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

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

	router.Register(app, userHandler)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
