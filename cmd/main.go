package main

// @title marketplace API
// @version 1.0
// @description API для тестового маркетплейса

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/vk_intern/docs"
	"github.com/vk_intern/internal/config"
	"github.com/vk_intern/internal/logger"
	"github.com/vk_intern/internal/repository"
	"github.com/vk_intern/routes"
)

func main() {
	cfg := config.MustLoad()

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	logger.Init("text")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// создание пула соединений к БД
	pool, err := repository.InitDB(ctx, cfg)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	defer pool.Close()

	// запуск миграций
	if err := repository.RunMigrations(cfg); err != nil {
		log.Fatal("Migration failed:", err)
	}

	routes.InitRoutes(app, cfg)
	log.Fatal(app.Listen(cfg.Server.Port))
}
