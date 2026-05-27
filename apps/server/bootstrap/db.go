package bootstrap

import (
	"fmt"
	"omni-pixel/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(env *Env) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Conversation{},
		&domain.Message{},
		&domain.Provider{},
		&domain.Model{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
