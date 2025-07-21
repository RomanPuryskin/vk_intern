package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vk_intern/internal/config"
	"github.com/vk_intern/internal/logger"
)

func RunMigrations(cfg *config.Config) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Storage.User,
			cfg.Storage.Password,
			cfg.Storage.Host,
			cfg.Storage.Port,
			cfg.Storage.Name,
		),
	)
	if err != nil {
		return fmt.Errorf("[RunMigrations| new migrate] %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[RunMigrations| up migrate] %w", err)
	}

	logger.L.Info("Applied migrations")
	return nil
}
