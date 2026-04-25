package main

import (
	"log"

	"omni-pixel/bootstrap"
)

func main() {
	env, err := bootstrap.NewEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	db, err := bootstrap.NewPostgresPool(env)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	app := bootstrap.NewApp(env, db)
	if err := app.Listen(env.ServerAddress); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
